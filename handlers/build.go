package handlers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/jonathanwthom/earl/models"
	"github.com/matoous/go-nanoid"
	"net/http"
	"strings"
)

var db *gorm.DB

type HandlerSet struct {
	GetLinksHandler      func(http.ResponseWriter, *http.Request)
	GetLinkHandler       func(http.ResponseWriter, *http.Request)
	CreateLinkHandler    func(http.ResponseWriter, *http.Request)
	CreateAccountHandler func(http.ResponseWriter, *http.Request)
}

func Build(database *gorm.DB) HandlerSet {
	db = database

	return HandlerSet{
		GetLinksHandler:      getLinksHandler,
		GetLinkHandler:       getLinkHandler,
		CreateLinkHandler:    createLinkHandler,
		CreateAccountHandler: createAccountHandler,
	}
}

func createLink(original string, req *http.Request) (*models.Link, error) {
	link := &models.Link{Original: original}

	// share header code
	auth := req.Header.Get("Authorization")
	if auth != "" {
		token := strings.ReplaceAll(auth, "basic ", "")
		account := &models.Account{}
		notFound := db.Where("token = ?", token).First(account).RecordNotFound()
		if notFound {
			return link, errors.New("No account with token")
		}

		link.AccountID = account.ID
	}

	err := link.Validate()
	if err != nil {
		return link, err
	}

	err = link.Shorten(req)
	if err != nil {
		return link, err
	}

	err = db.Create(link).Error
	if err != nil {
		return link, err
	}

	return link, nil
}

// need to return json
func createAccountHandler(w http.ResponseWriter, req *http.Request) {
	id, err := gonanoid.Nanoid()
	token := base64.StdEncoding.EncodeToString([]byte(id))
	if err != nil {
		http.Error(w, "Unable to create account", http.StatusInternalServerError)
		return
	}
	account := &models.Account{Token: token}

	err = db.Create(account).Error
	if err != nil {
		http.Error(w, "Unable to create account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Your account token is: %s\n", account.Token)
}
func getLinksHandler(w http.ResponseWriter, req *http.Request) {
	// share header fetch code
	auth := req.Header.Get("Authorization")
	account := &models.Account{}
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

	links := []models.Link{}
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

// return json
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
	fmt.Fprintf(w, msg, link.Shortened)
}

func getLinkHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	identifier := vars["identifier"]
	link := models.Link{Identifier: identifier}
	notFound := db.Where("identifier = ?", identifier).First(&link).RecordNotFound()

	if notFound == true {
		msg := fmt.Sprintf("Unable to find %s\n", identifier)
		http.Error(w, msg, http.StatusNotFound)
		return
	}

	url := link.Original

	// @todo: More logging
	// could log things about remote ip with https://godoc.org/github.com/oschwald/geoip2-golang
	view := &models.View{LinkID: link.ID, RemoteAddr: req.RemoteAddr}
	// handle errors
	db.Create(view)

	http.Redirect(w, req, url, 302)
}
