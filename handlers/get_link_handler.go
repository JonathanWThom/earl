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

	location, err := getLocationFromIP(req)
	if err != nil {
		log.Print(err)
	}

	view := &models.View{
		LinkID:     link.ID,
		RemoteAddr: req.RemoteAddr,
		UserAgent:  req.UserAgent(),
		Referer:    req.Referer(),
		Country:    location.CountryName,
		City:       location.CityName,
		ZipCode:    location.ZipCode,
	}
	err = db.Create(view).Error
	if err != nil {
		log.Print(err)
		http.Error(w, "Unable to redirect to link", http.StatusInternalServerError)
	}

	http.Redirect(w, req, url, 302)
}

func getLocationFromIP(req *http.Request) (models.Location, error) {
	rawIP := getIP(req)
	location := models.Location{}
	dec, err := ip2int(rawIP)
	if err != nil {
		return location, err

	}
	err = db.Where("ip_from <= ? and ip_to >= ?", dec, dec).First(&location).Error

	return location, err
}

func ip2int(raw string) (uint32, error) {
	ip, _, err := net.SplitHostPort(raw)
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint32([]byte(ip)), nil
}

func getIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}
