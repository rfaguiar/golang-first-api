package api

import (
	"github.com/gorilla/mux"
	"github.com/rfaguiar/golang-first-api/api/controller"
	"github.com/rfaguiar/golang-first-api/api/model"
	"log"
	"net/http"
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
	router.HandleFunc("/api-v1/user", controller.GetUsers).Methods("GET")
	router.HandleFunc("/api-v1/user/{id}", controller.GetAnUser).Methods("GET")
	router.HandleFunc("/api-v1/user", controller.CreateUser).Methods("POST")
	router.HandleFunc("/api-v1/user/{id}", controller.DeleteUser).Methods("DELETE")
	router.HandleFunc("/api-v1/user/{id}", controller.UpdateUser).Methods("PUT")
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
