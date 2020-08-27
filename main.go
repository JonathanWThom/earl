package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/jonathanwthom/earl/models"
	"github.com/jonathanwthom/earl/routes"
	"log"
	"os"
)

// @todo: Add caching
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

	routes.Serve(db)
}
