package handlers

import (
	"encoding/json"
	"errors"
	"github.com/jonathanwthom/earl/models"
	"log"
	"net/http"
)

func createLinkHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	url := req.FormValue("url")
	if url == "" {
		http.Error(w, "Missing parameter: url", http.StatusBadRequest)
		return
	}

	link, err := createLink(url, req)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	js, err := json.Marshal(link)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func createLink(original string, req *http.Request) (*models.Link, error) {
	link := &models.Link{Original: original}

	account, err := getAccountFromToken(req)
	if err != nil {
		return link, err
	}
	link.AccountID = account.ID

	err = link.Validate()
	if err != nil {
		log.Print(err)
		return link, errors.New("Invalid parameter: url")
	}

	err = link.Shorten(req)
	if err != nil {
		log.Print(err)
		return link, errors.New("Unable to create link")
	}

	err = db.Create(link).Error
	if err != nil {
		log.Print(err)
		return link, errors.New("Unable to create link")
	}

	return link, nil
}
