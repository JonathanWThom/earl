package main

import (
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
)

// @todo: Add caching
// @todo: Associate links with users/accounts (conditionally), add auth
// @todo: Add metrics storage and viewing
// @todo: Tests
// @todo: Break up main.go
// @todo: Add linter tool

var db *gorm.DB

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
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

	r := mux.NewRouter()

	// @todo: Add generic logging middleware around routes
	r.HandleFunc("/{identifier}", getLinkHandler).Methods("GET")
	r.HandleFunc("/links", createLinkHandler).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	http.ListenAndServe(":"+port, r)
}

type Link struct {
	gorm.Model
	Original   string `gorm:"not null"`
	Identifier string `gorm:"unique;not null"`
}

func createLink(original string) (*Link, error) {
	link := &Link{Original: original}
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

	link, err := createLink(url)
	if err != nil {
		http.Error(w, "Invalid parameter: url", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Your short url: %s\n", link.shortened(req))
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

	// @todo: Log things for user that the link belongs to

	http.Redirect(w, req, url, 302)
}
