package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jonathanwthom/earl/models"
	"net/http"
)

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
	// Could this also be done concurrently with redirect?
	view := &models.View{LinkID: link.ID, RemoteAddr: req.RemoteAddr}
	err := db.Create(view).Error
	if err != nil {
		http.Error(w, "Unable to redirect to link", http.StatusInternalServerError)
	}

	http.Redirect(w, req, url, 302)
}
