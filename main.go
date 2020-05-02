package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/health", healthCheck)
	//show log server address
	log.Print("Server listen http://localhost:9000")
	//UP server and listen http port 9000 using default http multiplexer, if error log message and kill api server
	log.Fatal(http.ListenAndServe(":9000", nil))
}

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
	log.Print("ALL / home")
	_, err := fmt.Fprint(responseWriter, "Server UP")
	if err != nil { // if error then log error and return status code 500 Internal Server Error
		log.Print(err.Error())
		responseWriter.WriteHeader(http.StatusInternalServerError)
	}
}
