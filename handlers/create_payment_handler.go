package handlers

import (
	"fmt"
	"github.com/jonathanwthom/earl/models"
	"github.com/jonathanwthom/earl/payments"
	"net/http"
	"strings"
)

func createPaymentHandler(w http.ResponseWriter, req *http.Request) {
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
		fmt.Fprintf(w, "Must pass basic Authorization header to make payment")
		return
	}

	payments.CreateCharge(account.Token)
}
