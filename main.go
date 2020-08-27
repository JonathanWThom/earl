package main

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/jonathanwthom/earl/handlers"
	"github.com/jonathanwthom/earl/models"
	"log"
	"net/http"
	"os"
)

// @todo: Add caching
// @todo: Break up main.go
// @todo: Add linter tool

var db *gorm.DB

func main() {
	godotenv.Load()
	var err error
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("No DATABASE_URL variable set")
	}
	db, err = gorm.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	db.LogMode(true)
	defer db.Close()
	db.AutoMigrate(&models.Link{})
	db.AutoMigrate(&models.Account{})
	db.AutoMigrate(&models.View{})

	r := mux.NewRouter()
	h := handlers.Build(db)

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
