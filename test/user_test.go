package controller_test

import (
	"net/http"
	"testing"
)

func TestEmptyTablePerson(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api-v1/user", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	if body := response.Body.String(); body != "[]\n" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}
