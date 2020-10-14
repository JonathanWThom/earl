package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/jonathanwthom/earl/models"
	"log"
	"net/http"
)

func getLinksHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println(req)

	account, err := getAccountFromToken(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, err.Error())
		return
	}

	links := []models.Link{}
	err = db.Where("account_id = ?", account.ID).Preload("Views").Find(&links).Error
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to fetch links")
		return
	}

	js, err := json.Marshal(links)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}
