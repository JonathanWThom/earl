package main

import (
	"github.com/joho/godotenv"
	"github.com/jonathanwthom/earl/database"
	"github.com/jonathanwthom/earl/routes"
)

// @todo: Add caching
// @todo: Add linter tool
// @todo: Add payment, account inspection
// @todo: Add tags/meta for links and queries based on tags

func main() {
	godotenv.Load()
	db := database.Init()
	defer db.Close()
	routes.Init(db)
}
