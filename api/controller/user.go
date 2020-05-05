package controller

import (
	"encoding/json"
	"github.com/rfaguiar/golang-first-api/api/model"
	"log"
	"net/http"
)

/*
	Show all users in a repository
*/
func GetUsers(responseWriter http.ResponseWriter, _ *http.Request) {
	log.Print("User controller: GET /api-v1/user")
	users := model.User{}.FindAll()
	err := json.NewEncoder(responseWriter).Encode(users)
	if err != nil {
		log.Print(err.Error())
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
}
