package models

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/matoous/go-nanoid"
	"net/http"
	"net/url"
)

type Link struct {
	gorm.Model `json:"-"`
	Original   string `gorm:"not null" json:"original""`
	Identifier string `gorm:"unique;not null"`
	Shortened  string `gorm:"unique" json:"shortened"`
	AccountID  uint   `json:"-"`
	Views      []View `json:"views"`
}

func (link *Link) Shorten(req *http.Request) error {
	identifier, err := gonanoid.ID(6)
	if err != nil {
		return err
	}

	link.Identifier = identifier
	link.Shortened = link.shortened(req)

	return nil
}

func (link *Link) Validate() error {
	original := link.Original
	_, err := url.ParseRequestURI(original)
	if err != nil {
		return errors.New("Invalid URL")
	}

	u, err := url.Parse(original)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return errors.New("Invalid URL")
	}

	// @todo: Could also try a GET request on link to make sure it exists

	return nil
}

func (link *Link) shortened(request *http.Request) string {
	return fmt.Sprintf("https://%s/%s", request.Host, link.Identifier)
}
