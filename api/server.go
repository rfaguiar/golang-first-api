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
	s.get("/", controller.Home)
	s.get("/health", controller.HealthCheck)
	s.get("/api-v1/user", controller.GetUsers)
	s.get("/api-v1/user/{id}", controller.GetAnUser)
	s.post("/api-v1/user", controller.CreateUser)
	s.del("/api-v1/user/{id}", controller.DeleteUser)
	s.put("/api-v1/user/{id}", controller.UpdateUser)
}

func (s *Server) get(path string, f func(http.ResponseWriter, *http.Request)) {
	s.Router.HandleFunc(path, f).Methods(http.MethodGet)
}

func (s *Server) post(path string, f func(http.ResponseWriter, *http.Request)) {
	s.Router.HandleFunc(path, f).Methods(http.MethodPost)
}

func (s *Server) del(path string, f func(http.ResponseWriter, *http.Request)) {
	s.Router.HandleFunc(path, f).Methods(http.MethodDelete)
}

func (s *Server) put(path string, f func(http.ResponseWriter, *http.Request)) {
	s.Router.HandleFunc(path, f).Methods(http.MethodPut)
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
