package api

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
	Startup server API
*/
func Run() {
	model.User{}.InitializeDatabase()
	router := mux.NewRouter()
	router.Use(jsonMiddleware)
	router.HandleFunc("/", home).Methods("GET")
	router.HandleFunc("/health", healthCheck).Methods("GET")
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
	index := -1
	for key, u := range model.UserRepo {
		if u.Id == id {
			index = key
		}
	}
	if index < 0 {
		responseWriter.WriteHeader(http.StatusNotFound)
		return
	}
	//update user in a repository
	var user model.User
	err = json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		log.Print(err.Error())
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	user.Id = id
	model.UserRepo[index] = user
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
	//find user by id
	index := -1
	for key, u := range model.UserRepo {
		if u.Id == id {
			index = key
		}
	}
	if index < 0 {
		responseWriter.WriteHeader(http.StatusNotFound)
		return
	}
	//delete user in a repository
	leftSlice := model.UserRepo[0:index]
	rightSlice := model.UserRepo[index+1:]
	model.UserRepo = append(leftSlice, rightSlice...)
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
	user.Id = len(model.UserRepo) + 1
	model.UserRepo = append(model.UserRepo, user)
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
	for _, user := range model.UserRepo {
		if user.Id == id {
			err = json.NewEncoder(responseWriter).Encode(user)
			if err != nil {
				log.Print(err.Error())
				responseWriter.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
	}
	responseWriter.WriteHeader(http.StatusNotFound)
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

/*
	show health check status UP and set status code 200 OK
	if error print log and set status 500 Internal Server Error
*/
func healthCheck(responseWriter http.ResponseWriter, _ *http.Request) {
	log.Print("GET /health")
	health := model.Health{Status: "UP"}
	err := json.NewEncoder(responseWriter).Encode(health)
	if err != nil {
		log.Print(err.Error())
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
}

/*
	show friendly message
*/
func home(responseWriter http.ResponseWriter, _ *http.Request) {
	responseWriter.Header().Set("Content-Type", "text/plain")
	log.Print("GET / home")
	_, err := fmt.Fprint(responseWriter, "Server UP")
	if err != nil { // if error then log error and return status code 500 Internal Server Error
		log.Print(err.Error())
		responseWriter.WriteHeader(http.StatusInternalServerError)
	}
}
