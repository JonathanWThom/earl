package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/jonathanwthom/earl/models"
	"net/http"
	"strings"
)

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
