package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/jonathanwthom/earl/models"
	"log"
	"os"
)

func Init() *gorm.DB {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("No DATABASE_URL variable set")
	}

	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	db.LogMode(true)

	db.AutoMigrate(&models.Link{})
	db.AutoMigrate(&models.Account{})
	db.AutoMigrate(&models.View{})

	return db
}
