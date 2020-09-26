package routes

import (
	middleware "github.com/gorilla/handlers"
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

	// Order matters
	r.Handle(
		"/",
		http.FileServer(http.Dir("./static")),
	)
	r.Handle(
		"/links", middleware.LoggingHandler(os.Stdout, http.HandlerFunc(h.GetLinksHandler)),
	).Methods("GET")
	r.Handle(
		"/{identifier}",
		middleware.LoggingHandler(os.Stdout, http.HandlerFunc(h.GetLinkHandler)),
	).Methods("GET")
	r.Handle(
		"/links",
		middleware.LoggingHandler(os.Stdout, http.HandlerFunc(h.CreateLinkHandler)),
	).Methods("POST")
	r.Handle(
		"/accounts",
		middleware.LoggingHandler(os.Stdout, http.HandlerFunc(h.CreateAccountHandler)),
	).Methods("POST")
	r.Handle(
		"/payments",
		middleware.LoggingHandler(os.Stdout, http.HandlerFunc(h.CreatePaymentHandler)),
	).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	http.ListenAndServe(":"+port, r)
}
