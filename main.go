package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

/*
	Type Health for use health checks endpoint
	Ex: status: "UP"
*/
type Health struct {
	Status string `json:"status"`
}

/*
	User type for use CAD
*/
type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

/*
	Repository for users type
*/
var userRepo []User

/*
	Startup server API
*/
func main() {
	initializeDatabase()
	router := mux.NewRouter()
	router.Use(jsonMiddleware)
	router.HandleFunc("/", home).Methods("GET")
	router.HandleFunc("/health", healthCheck).Methods("GET")
	router.HandleFunc("/api-v1/user", getUsers).Methods("GET")
	router.HandleFunc("/api-v1/user/{id}", getAnUser).Methods("GET")
	//show log server address
	log.Print("Server listen http://localhost:9000")
	//UP server and listen http port 9000 using default http multiplexer, if error log message and kill api server
	log.Fatal(http.ListenAndServe(":9000", router))
}

/*
	Initialize in memory database
*/
func initializeDatabase() {
	userRepo = append(userRepo,
		User{Id: 1, Name: "Jorge", Age: 20},
		User{Id: 2, Name: "Jhon", Age: 33})
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
	for _, user := range userRepo {
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
	err := json.NewEncoder(responseWriter).Encode(userRepo)
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
	health := Health{"UP"}
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
