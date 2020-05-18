package controller

import (
	"encoding/json"
	"fmt"
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

	if err := json.NewEncoder(responseWriter).Encode(users); err != nil {
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
	log.Printf("User controller: GET /api-v1/user/%v", id)
	user := model.User{}.FindById(id)
	if user == nil {
		responseWriter.WriteHeader(http.StatusNotFound)
		return
	}
	if err = json.NewEncoder(responseWriter).Encode(user); err != nil {
		log.Print(err.Error())
		responseWriter.WriteHeader(http.StatusInternalServerError)
	}
}

/*
	Create new user and save in a repository and set location
*/
func CreateUser(responseWriter http.ResponseWriter, request *http.Request) {
	log.Print("User controller: POST /api-v1/user")
	var user model.User
	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		log.Print(err.Error())
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := user.Save(); err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(responseWriter, "{error: %v}", err.Error())
		return
	}
	responseWriter.Header().Set("location", fmt.Sprintf("/api-v1/user/%v", user.Id))
	responseWriter.WriteHeader(http.StatusCreated)
}

/*
	Delete user by param id
*/
func DeleteUser(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Print(err.Error())
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("User controller: DELETE /api-v1/user/%v", id)
	user := model.User{}.FindById(id)
	if user == nil {
		responseWriter.WriteHeader(http.StatusNotFound)
		return
	}
	if err := user.Remove(); err != nil {
		_, _ = fmt.Fprintf(responseWriter, "{error: %v}", err.Error())
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
}

/*
	Update user using parameter id and attributes inside body
*/
func UpdateUser(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Print(err.Error())
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("PUT /api-v1/user/%v", id)
	//find user by id
	user := model.User{}.FindById(id)
	if user == nil {
		responseWriter.WriteHeader(http.StatusNotFound)
		return
	}
	var userToSave model.User
	err = json.NewDecoder(request.Body).Decode(&userToSave)
	if err != nil {
		log.Print(err.Error())
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := userToSave.Update(user.Id); err != nil {
		_, _ = fmt.Fprintf(responseWriter, "{error: %v}", err.Error())
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
}
