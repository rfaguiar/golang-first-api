package controller_test

import (
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
	//create users
	db := database.Current()
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("INSERT INTO person(name, age) VALUES($1, $2)")
	stmt.Exec("Jorge", 20)
	stmt.Exec("Jhon", 33)
	tx.Commit()
	stmt.Close()

	//get /users
	req, _ := http.NewRequest("GET", "/api-v1/user", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	expected := "[{\"id\":1,\"name\":\"Jorge\",\"age\":20},{\"id\":2,\"name\":\"Jhon\",\"age\":33}]\n"
	if body := response.Body.String(); body != expected {
		t.Errorf("Expected an users array. Got %s", body)
	}

	//remove users
	txr, _ := db.Begin()
	stmts, _ := txr.Prepare("DELETE FROM person")
	stmts.Exec()
	txr.Commit()
	stmts.Close()
}
