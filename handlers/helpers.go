package handlers

import (
	"errors"
	"github.com/jonathanwthom/earl/models"
	"net/http"
	"strings"
)

func getAccountFromToken(req *http.Request) (*models.Account, error) {
	account := &models.Account{}
	auth := req.Header.Get("Authorization")

	if auth != "" {
		token := strings.ReplaceAll(auth, "basic ", "")
		notFound := db.Where("token = ?", token).First(account).RecordNotFound()
		if notFound {
			return account, errors.New("No account with token")
		}

		return account, nil
	}

	return account, errors.New("Missing Authorization header")
}
