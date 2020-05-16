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
	stmt, err := tx.Prepare("INSERT INTO person(name, age) VALUES($1, $2)")
	if err != nil {
		t.Error(err)
	}
	_, err = stmt.Exec("Jorge", 20)
	if err != nil {
		t.Error(err)
	}
	_, err = stmt.Exec("Jhon", 33)
	if err != nil {
		t.Error(err)
	}
	if err != nil {
		t.Error(err)
	}
	if err := tx.Commit(); err != nil {
		t.Error(err)
	}
	if err := stmt.Close(); err != nil {
		t.Error(err)
	}

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
	stmts, err := txr.Prepare("DELETE FROM person")
	if err != nil {
		t.Error(err)
	}
	_, err = stmts.Exec()
	if err != nil {
		t.Error(err)
	}
	if err := txr.Commit(); err != nil {
		t.Error(err)
	}
	if err := stmts.Close(); err != nil {
		t.Error(err)
	}
}
