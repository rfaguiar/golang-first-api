package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
)

const (
	host     = "host"
	port     = "port"
	user     = "user"
	password = "password"
	dbname   = "dbname"
)

var db *sql.DB

func Open() *sql.DB {
	fmt.Println("Connecting into postgress database")

	dbInfo := getDbInfo()
	portInfo, err := strconv.Atoi(dbInfo[port])
	if err != nil {
		log.Fatalf(err.Error())
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbInfo[host], portInfo, dbInfo[user], dbInfo[password], dbInfo[dbname])

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("could not connect database %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("could not connect database %v", err)
	}
	fmt.Println("Successfully connected")
	return db
}

func Close() {
	if err := db.Close(); err != nil {
		log.Fatalf("could not close database %v", err)
	}
}

func getDbInfo() map[string]string {
	result := make(map[string]string, 0)
	envHost, ok := os.LookupEnv("ENV_DB_HOST")
	if !ok {
		log.Fatal("ENV_DB_HOST not found")
	}
	result[host] = envHost

	envPort, ok := os.LookupEnv("ENV_DB_PORT")
	if !ok {
		log.Fatal("ENV_DB_PORT not found")
	}
	result[port] = envPort

	envUser, ok := os.LookupEnv("ENV_DB_USER")
	if !ok {
		log.Fatal("ENV_DB_USER not found")
	}
	result[user] = envUser

	envPass, ok := os.LookupEnv("ENV_DB_PASSWORD")
	if !ok {
		log.Fatal("ENV_DB_PASSWORD not found")
	}
	result[password] = envPass

	envDbname, ok := os.LookupEnv("ENV_DB_DBNAME")
	if !ok {
		log.Fatal("ENV_DB_DBNAME not found")
	}
	result[dbname] = envDbname

	return result
}
