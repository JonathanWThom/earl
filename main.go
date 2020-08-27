package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/matoous/go-nanoid"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// @todo: Add caching
// @todo: Associate links with users/accounts (conditionally), add auth
// @todo: Add metrics storage and viewing
// @todo: Tests
// @todo: Break up main.go
// @todo: Add linter tool

var db *gorm.DB

func main() {
	godotenv.Load()
	var err error
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("No DATABASE_URL variable set")
	}
	db, err = gorm.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	db.LogMode(true)
	defer db.Close()
	db.AutoMigrate(&Link{})
	db.AutoMigrate(&Account{})
	db.AutoMigrate(&View{})

	r := mux.NewRouter()

	// @todo: Add generic logging middleware around routes
	// Order matters
	r.HandleFunc("/links", getLinksHandler).Methods("GET")
	r.HandleFunc("/{identifier}", getLinkHandler).Methods("GET")
	r.HandleFunc("/links", createLinkHandler).Methods("POST")
	r.HandleFunc("/accounts", createAccountHandler).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	http.ListenAndServe(":"+port, r)
}

type Link struct {
	gorm.Model `json:"-"`
	Original   string `gorm:"not null" json:"original""`
	Identifier string `gorm:"unique;not null" json:"identifier"`
	AccountID  uint   `json:"-"`
	Views      []View `json:"views"`
}

type View struct {
	ID         uint       `gorm:"primary_key" json:"-"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"-"`
	DeletedAt  *time.Time `sql:"index" json:"-"`
	LinkID     uint       `json:"-"`
	RemoteAddr string     `json:"remoteAddr"`
}

type Account struct {
	gorm.Model
	Token string `gorm:"unique;not null"`
	Links []Link
}

func createAccountHandler(w http.ResponseWriter, req *http.Request) {
	id, err := gonanoid.Nanoid()
	token := base64.StdEncoding.EncodeToString([]byte(id))
	if err != nil {
		http.Error(w, "Unable to create account", http.StatusInternalServerError)
		return
	}
	account := &Account{Token: token}

	err = db.Create(account).Error
	if err != nil {
		http.Error(w, "Unable to create account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Your account token: %s. Pass token as Authorization header: `basic your-token-goes-here`\n", account.Token)
}

func getLinksHandler(w http.ResponseWriter, req *http.Request) {
	auth := req.Header.Get("Authorization")
	account := &Account{}
	if auth != "" {
		token := strings.ReplaceAll(auth, "basic ", "")
		notFound := db.Where("token = ?", token).First(account).RecordNotFound()
		if notFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Unable to find account")
			return
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Must pass basic Authorization header to read links")
		return
	}

	links := []Link{}
	err := db.Where("account_id = ?", account.ID).Preload("Views").Find(&links).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to fetch links")
		return
	}

	js, err := json.Marshal(links)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func createLink(original string, req *http.Request) (*Link, error) {
	link := &Link{Original: original}

	auth := req.Header.Get("Authorization")
	if auth != "" {
		token := strings.ReplaceAll(auth, "basic ", "")
		account := &Account{}
		notFound := db.Where("token = ?", token).First(account).RecordNotFound()
		if notFound {
			return link, errors.New("No account with token")
		}

		link.AccountID = account.ID
	}

	err := link.validate()
	if err != nil {
		return link, err
	}

	err = link.shorten()
	if err != nil {
		return link, err
	}

	err = db.Create(link).Error
	if err != nil {
		return link, err
	}

	return link, nil
}

func (link *Link) shorten() error {
	identifier, err := gonanoid.ID(6)
	if err != nil {
		return err
	}

	link.Identifier = identifier

	return nil
}

func (link *Link) validate() error {
	original := link.Original
	_, err := url.ParseRequestURI(original)
	if err != nil {
		return errors.New("Invalid URL")
	}

	u, err := url.Parse(original)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return errors.New("Invalid URL")
	}

	// @todo: Could also try a GET request on link to make sure it exists

	return nil
}

func (link *Link) shortened(request *http.Request) string {
	return fmt.Sprintf("%s/%s", request.Host, link.Identifier)
}

func createLinkHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	url := req.FormValue("url")
	if url == "" {
		http.Error(w, "Missing parameter: url", http.StatusBadRequest)
		return
	}

	link, err := createLink(url, req)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid parameter: url", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	msg := "Your short url (created without account): %s\n"
	if link.AccountID != 0 {
		msg = "Your short url (created for account): %s\n"
	}
	fmt.Fprintf(w, msg, link.shortened(req))
}

func getLinkHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	identifier := vars["identifier"]
	link := Link{Identifier: identifier}
	notFound := db.Where("identifier = ?", identifier).First(&link).RecordNotFound()

	if notFound == true {
		msg := fmt.Sprintf("Unable to find %s\n", identifier)
		http.Error(w, msg, http.StatusNotFound)
		return
	}

	url := link.Original

	// @todo: More logging
	// could log things about remote ip with https://godoc.org/github.com/oschwald/geoip2-golang
	view := &View{LinkID: link.ID, RemoteAddr: req.RemoteAddr}
	// handle errors
	db.Create(view)

	http.Redirect(w, req, url, 302)
}
