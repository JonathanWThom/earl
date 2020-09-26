package handlers

import (
	"github.com/jinzhu/gorm"
	"net/http"
)

var db *gorm.DB

type HandlerSet struct {
	GetLinksHandler      func(http.ResponseWriter, *http.Request)
	GetLinkHandler       func(http.ResponseWriter, *http.Request)
	CreateLinkHandler    func(http.ResponseWriter, *http.Request)
	CreateAccountHandler func(http.ResponseWriter, *http.Request)
	CreatePaymentHandler func(http.ResponseWriter, *http.Request)
}

func Init(database *gorm.DB) HandlerSet {
	db = database

	return HandlerSet{
		GetLinksHandler:      getLinksHandler,
		GetLinkHandler:       getLinkHandler,
		CreateLinkHandler:    createLinkHandler,
		CreateAccountHandler: createAccountHandler,
		CreatePaymentHandler: createPaymentHandler,
	}
}
