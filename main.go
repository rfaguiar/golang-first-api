package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(responseWriter http.ResponseWriter, request *http.Request) {
		//show friendly message
		_, err := fmt.Fprint(responseWriter, "Server UP")
		if err != nil { // if error then log error and return status code 500 Internal Server Error
			log.Print(err.Error())
			responseWriter.WriteHeader(http.StatusInternalServerError)
		}
	})
	//show log server address
	log.Print("Server listen http://localhost:9000")
	//UP server and listen http port 9000 using default http multiplexer, if error log message and kill api server
	log.Fatal(http.ListenAndServe(":9000", nil))
}
