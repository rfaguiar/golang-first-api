package controller

import (
	"encoding/json"
	"github.com/rfaguiar/golang-first-api/api/model"
	"log"
	"net/http"
)

/*
	show health check status UP and set status code 200 OK
	if error print log and set status 500 Internal Server Error
*/
func HealthCheck(responseWriter http.ResponseWriter, _ *http.Request) {
	log.Print("Health controller: GET /health")
	health := model.Health{Status: "UP"}
	err := json.NewEncoder(responseWriter).Encode(health)
	if err != nil {
		log.Print(err.Error())
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
}
