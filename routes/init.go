package routes

import (
	middleware "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/jonathanwthom/earl/handlers"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
)

func Init(db *gorm.DB) {
	r := mux.NewRouter()
	h := handlers.Init(db)

	// Order matters
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

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"Authorization"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	handler := c.Handler(r)
	http.ListenAndServe(":"+port, handler)
}
