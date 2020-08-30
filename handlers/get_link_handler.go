package handlers

import (
	"encoding/binary"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jonathanwthom/earl/models"
	"log"
	"net"
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

	location, err := getLocationFromIP(req.RemoteAddr)
	if err != nil {
		log.Print(err)
	}
	log.Print(location)

	view := &models.View{
		LinkID:     link.ID,
		RemoteAddr: req.RemoteAddr,
		UserAgent:  req.UserAgent(),
		Referer:    req.Referer(),
		Location:   location,
	}
	err = db.Create(view).Error
	if err != nil {
		log.Print(err)
		http.Error(w, "Unable to redirect to link", http.StatusInternalServerError)
	}

	http.Redirect(w, req, url, 302)
}

func getLocationFromIP(rawIP string) (models.Location, error) {
	//dec := ip2int(rawIP)
	location := models.Location{}
	//err := db.Where("ip_from <= ? and ip_to >= ?", dec, dec).First(&location).Error

	return location, nil
}

func ip2int(raw string) uint32 {
	fmt.Println(raw)
	ip := net.ParseIP(raw)
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}
