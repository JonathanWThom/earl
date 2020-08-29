package handlers

import (
	"encoding/json"
	"errors"
	"github.com/jonathanwthom/earl/models"
	"net/http"
	"strings"
)

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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	js, err := json.Marshal(link)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(js)
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
		return link, errors.New("Invalid parameter: url")
	}

	err = link.Shorten(req)
	if err != nil {
		return link, errors.New("Unable to create link")
	}

	err = db.Create(link).Error
	if err != nil {
		return link, errors.New("Unable to create link")
	}

	return link, nil
}
