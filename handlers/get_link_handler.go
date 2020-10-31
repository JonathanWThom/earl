package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jonathanwthom/earl/models"
	"log"
	"math/big"
	"net"
	"net/http"
	"net/url"
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

	original := link.Original

	// @todo: Break this out
	go func() {
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
		}
	}()

	params := req.URL.Query()
	parsed, _ := url.Parse(original)
	q := parsed.Query()
	for key, value := range params {
		q.Set(key, value[0])
	}
	parsed.RawQuery = q.Encode()

	http.Redirect(w, req, parsed.String(), 302)
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
