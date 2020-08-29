package handlers

import (
	"encoding/base64"
	"encoding/json"
	"github.com/jonathanwthom/earl/models"
	"github.com/matoous/go-nanoid"
	"net/http"
)

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

	js, err := json.Marshal(account)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(js)
}
