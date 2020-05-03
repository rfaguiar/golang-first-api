package model

/*
	Type Health for use health checks endpoint
	Ex: status: "UP"
*/
type Health struct {
	Status string `json:"status"`
}
