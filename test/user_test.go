package controller_test

import (
	"fmt"
	"github.com/rfaguiar/golang-first-api/database"
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

func TestFindAllPerson(t *testing.T) {
	createPerson()

	//get /users
	req, _ := http.NewRequest("GET", "/api-v1/user", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	expected := "[{\"id\":1,\"name\":\"Jorge\",\"age\":20},{\"id\":2,\"name\":\"Jhon\",\"age\":33}]\n"
	if body := response.Body.String(); body != expected {
		t.Errorf("Expected an users array. Got %s", body)
	}

	removePerson()
}

func TestFindAnExistsPerson(t *testing.T) {
	ids := createPerson()

	//get /users/1
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api-v1/user/%d", ids[0]), nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	expected := fmt.Sprintf("{\"id\":%d,\"name\":\"Jorge\",\"age\":20}\n", ids[0])
	if body := response.Body.String(); body != expected {
		t.Errorf("Expected an user but. Got %s", body)
	}

	removePerson()
}

func TestFindAnNotExistsPerson(t *testing.T) {
	//get /users/0
	req, _ := http.NewRequest("GET", "/api-v1/user/0", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
	expected := ""
	if body := response.Body.String(); body != expected {
		t.Errorf("Expected empty but. Got %s", body)
	}
}

func createPerson() []int {
	db := database.Current()
	result := make([]int, 0)
	var personId int
	db.QueryRow("INSERT INTO person(name, age) VALUES($1, $2) RETURNING id", "Jorge", 20).Scan(&personId)
	result = append(result, personId)
	db.QueryRow("INSERT INTO person(name, age) VALUES($1, $2) RETURNING id", "Jhon", 33).Scan(&personId)
	result = append(result, personId)
	return result
}

func removePerson() {
	db := database.Current()
	txr, _ := db.Begin()
	stmts, _ := txr.Prepare("DELETE FROM person")
	stmts.Exec()
	txr.Commit()
	stmts.Close()
}
