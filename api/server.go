package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rfaguiar/golang-first-api/api/controller"
	"github.com/rfaguiar/golang-first-api/api/model"
	"log"
	"net/http"
	"strconv"
)

/*
	Startup server API
*/
func Run() {
	model.User{}.InitializeDatabase()
	router := mux.NewRouter()
	router.Use(jsonMiddleware)
	router.HandleFunc("/", controller.Home).Methods("GET")
	router.HandleFunc("/health", controller.HealthCheck).Methods("GET")
	router.HandleFunc("/api-v1/user", getUsers).Methods("GET")
	router.HandleFunc("/api-v1/user/{id}", getAnUser).Methods("GET")
	router.HandleFunc("/api-v1/user", createUser).Methods("POST")
	router.HandleFunc("/api-v1/user/{id}", deleteUser).Methods("DELETE")
	router.HandleFunc("/api-v1/user/{id}", updateUser).Methods("PUT")
	//show log server address
	log.Print("Server listen http://localhost:9000")
	//UP server and listen http port 9000 using default http multiplexer, if error log message and kill api server
	log.Fatal(http.ListenAndServe(":9000", router))

}

/*
	Middleware for set all requests content type application json
*/
func jsonMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		handler.ServeHTTP(responseWriter, request)
	})
}

/*
	Update user using parameter id and attributes inside body
*/
func updateUser(responseWriter http.ResponseWriter, request *http.Request) {
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
	userToSave.Update(user.Id)
}

/*
	Delete user by param id
*/
func deleteUser(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Print(err.Error())
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("DELETE /api-v1/user/%v", id)
	user := model.User{}.FindById(id)
	if user == nil {
		responseWriter.WriteHeader(http.StatusNotFound)
		return
	}
	user.Remove()
}

/*
	Create new user and save in a repository and set location
*/
func createUser(responseWriter http.ResponseWriter, request *http.Request) {
	log.Print("POST /api-v1/user")
	var user model.User
	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		log.Print(err.Error())
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	user.Save()
	responseWriter.Header().Set("location", fmt.Sprintf("/api-v1/user/%v", user.Id))
	responseWriter.WriteHeader(http.StatusCreated)
}

/*
	show an user
*/
func getAnUser(responseWriter http.ResponseWriter, request *http.Request) {
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

/*
	Show all users in a repository
*/
func getUsers(responseWriter http.ResponseWriter, _ *http.Request) {
	log.Print("GET /api-v1/user")
	users := model.User{}.FindAll()
	err := json.NewEncoder(responseWriter).Encode(users)
	if err != nil {
		log.Print(err.Error())
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
}
