package api

import (
	"github.com/gorilla/mux"
	"github.com/rfaguiar/golang-first-api/api/controller"
	"log"
	"net/http"
)

type Server struct {
	Router *mux.Router
}

/*
	Startup server API
*/
func (s *Server) Run() {
	s.Initialize()
	log.Print("Server startup and listen http://localhost:9000")
	log.Fatal(http.ListenAndServe(":9000", s.Router))
}

/*
	Create new router and initialize http routes with controllers handlers
*/
func (s *Server) Initialize() {
	s.Router = mux.NewRouter()
	s.initializeRoutes()
}

func (s *Server) initializeRoutes() {
	s.Router.Use(jsonMiddleware)
	s.Router.HandleFunc("/", controller.Home).Methods("GET")
	s.Router.HandleFunc("/health", controller.HealthCheck).Methods("GET")
	s.Router.HandleFunc("/api-v1/user", controller.GetUsers).Methods("GET")
	s.Router.HandleFunc("/api-v1/user/{id}", controller.GetAnUser).Methods("GET")
	s.Router.HandleFunc("/api-v1/user", controller.CreateUser).Methods("POST")
	s.Router.HandleFunc("/api-v1/user/{id}", controller.DeleteUser).Methods("DELETE")
	s.Router.HandleFunc("/api-v1/user/{id}", controller.UpdateUser).Methods("PUT")
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
