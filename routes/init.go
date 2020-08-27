package routes

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/jonathanwthom/earl/handlers"
	"log"
	"net/http"
	"os"
)

func Init(db *gorm.DB) {
	r := mux.NewRouter()
	h := handlers.Init(db)

	// @todo: Add generic logging middleware around routes
	// Order matters
	r.HandleFunc("/links", h.GetLinksHandler).Methods("GET")
	r.HandleFunc("/{identifier}", h.GetLinkHandler).Methods("GET")
	r.HandleFunc("/links", h.CreateLinkHandler).Methods("POST")
	r.HandleFunc("/accounts", h.CreateAccountHandler).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	http.ListenAndServe(":"+port, r)
}
