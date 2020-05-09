package main

import (
	"github.com/rfaguiar/golang-first-api/api"
	"github.com/rfaguiar/golang-first-api/database"
)

func main() {
	database.ExecuteMigration()
	api.Run()
}
