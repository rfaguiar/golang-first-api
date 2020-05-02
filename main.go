package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

/*
	Startup server API
*/
func main() {
	router := mux.NewRouter()
	router.Use(jsonMiddleware)
	router.HandleFunc("/", home).Methods("GET")
	router.HandleFunc("/health", healthCheck).Methods("GET")
	router.HandleFunc("/api-v1/user", getUsers).Methods("GET")
	//show log server address
	log.Print("Server listen http://localhost:9000")
	//UP server and listen http port 9000 using default http multiplexer, if error log message and kill api server
	log.Fatal(http.ListenAndServe(":9000", router))
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
	Show all users in a repository
*/
func getUsers(responseWriter http.ResponseWriter, request *http.Request) {
	log.Print("GET /user")
	users := []User{
		{Id: 1, Name: "Jorge", Age: 20},
		{Id: 2, Name: "Jhon", Age: 33}}
	js, err := json.Marshal(users)
	if err != nil {
		log.Print(err.Error())
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = responseWriter.Write(js)
	if err != nil {
		log.Print(err.Error())
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
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
	Type Health for use health checks endpoint
	Ex: status: "UP"
*/
type Health struct {
	Status string `json:"status"`
}

/*
	show health check status UP and set status code 200 OK
	if error print log and set status 500 Internal Server Error
*/
func healthCheck(responseWriter http.ResponseWriter, _ *http.Request) {
	log.Print("ALL /health")
	h := Health{"UP"}
	js, er := json.Marshal(h)
	if er != nil {
		log.Print(er.Error())
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	responseWriter.WriteHeader(http.StatusOK)
	_, err := responseWriter.Write(js)
	if err != nil {
		log.Print(err.Error())
		responseWriter.WriteHeader(http.StatusInternalServerError)
	}
}

/*
	show friendly message
*/
func home(responseWriter http.ResponseWriter, _ *http.Request) {
	responseWriter.Header().Set("Content-Type", "text/plain")
	log.Print("ALL / home")
	_, err := fmt.Fprint(responseWriter, "Server UP")
	if err != nil { // if error then log error and return status code 500 Internal Server Error
		log.Print(err.Error())
		responseWriter.WriteHeader(http.StatusInternalServerError)
	}
}
