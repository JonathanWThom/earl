package handlers

import (
	"encoding/base64"
	"fmt"
	"github.com/jonathanwthom/earl/models"
	"github.com/matoous/go-nanoid"
	"net/http"
)

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
