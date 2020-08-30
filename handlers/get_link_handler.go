package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jonathanwthom/earl/models"
	"log"
	"math/big"
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
		LinkID:    link.ID,
		UserAgent: req.UserAgent(),
		Referer:   req.Referer(),
		Country:   location.CountryName,
		City:      location.CityName,
		ZipCode:   location.ZipCode,
	}
	err = db.Create(view).Error
	if err != nil {
		log.Print(err)
		http.Error(w, "Unable to redirect to link", http.StatusInternalServerError)
	}

	http.Redirect(w, req, url, 302)
}

func getLocationFromIP(req *http.Request) (models.Location, error) {
	// This is need for Heroku, since it uses proxies
	ip := req.Header.Get("X-FORWARDED-FOR")
	dec := IP4toInt(net.ParseIP(ip))
	location := models.Location{}
	err := db.Where("ip_from <= ? and ip_to >= ?", dec, dec).First(&location).Error

	return location, err
}

func IP4toInt(IPv4Address net.IP) int64 {
	IPv4Int := big.NewInt(0)
	IPv4Int.SetBytes(IPv4Address.To4())
	return IPv4Int.Int64()
}
