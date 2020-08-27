package handlers

import (
	"errors"
	"fmt"
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
