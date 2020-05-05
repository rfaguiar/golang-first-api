package controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rfaguiar/golang-first-api/api/model"
	"log"
	"net/http"
	"strconv"
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

/*
	show an user
*/
func GetAnUser(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Print(err.Error())
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("GET /api-v1/user/%v", id)
	user := model.User{}.FindById(id)
	if user == nil {
		responseWriter.WriteHeader(http.StatusNotFound)
		return
	}
	err = json.NewEncoder(responseWriter).Encode(user)
	if err != nil {
		log.Print(err.Error())
		responseWriter.WriteHeader(http.StatusInternalServerError)
	}
}
