package main

import (
	"github.com/rfaguiar/golang-first-api/api"
	"github.com/rfaguiar/golang-first-api/database"
)

func main() {
	db := database.Open()
	defer db.Close()
	database.ExecuteMigration(db)
	server := api.Server{}
	server.Run()
}
