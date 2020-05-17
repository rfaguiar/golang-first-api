package controller_test

import (
	"fmt"
	"github.com/rfaguiar/golang-first-api/database"
	"net/http"
	"strings"
	"testing"
)

func TestEmptyTablePerson(t *testing.T) {
	response := getRequest(t, "/api-v1/user")
	checkResponseCode(t, http.StatusOK, response.Code)
	if body := response.Body.String(); body != "[]\n" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestFindAllPerson(t *testing.T) {
	createPerson()

	//get /users
	response := getRequest(t, "/api-v1/user")
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
	response := getRequest(t, fmt.Sprintf("/api-v1/user/%d", ids[0]))
	checkResponseCode(t, http.StatusOK, response.Code)
	expected := fmt.Sprintf("{\"id\":%d,\"name\":\"Jorge\",\"age\":20}\n", ids[0])
	if body := response.Body.String(); body != expected {
		t.Errorf("Expected an user but. Got %s", body)
	}

	removePerson()
}

func TestFindAnNotExistsPerson(t *testing.T) {
	//get /users/0
	response := getRequest(t, "/api-v1/user/0")
	checkResponseCode(t, http.StatusNotFound, response.Code)
	expected := ""
	if body := response.Body.String(); body != expected {
		t.Errorf("Expected empty but. Got %s", body)
	}
}

func TestFindAnIncorrectFormatPersonId(t *testing.T) {
	//get /users/0
	response := getRequest(t, "/api-v1/user/a")
	checkResponseCode(t, http.StatusInternalServerError, response.Code)
	expected := ""
	if body := response.Body.String(); body != expected {
		t.Errorf("Expected empty but. Got %s", body)
	}
}

func TestCreatePerson(t *testing.T) {
	//post /users
	user := "{\"name\":\"Ana\",\"age\":23}\n"
	response := postRequest(t, "/api-v1/user", strings.NewReader(user))
	checkResponseCode(t, http.StatusCreated, response.Code)
	expected := ""
	if body := response.Body.String(); body != expected {
		t.Errorf("Expected empty but. Got %s", body)
	}

	var id int
	db := database.Current()
	row := db.QueryRow("select id from person where name = $1 and age = $2", "Ana", "23")
	if err := row.Scan(&id); err != nil {
		t.Errorf("Error when find user by name and age: %s", err.Error())
	}

	locationExp := fmt.Sprintf("/api-v1/user/%d", id)
	if location := response.Header().Get("location"); location != locationExp {
		t.Errorf("Expected location %s but. Got %s", locationExp, location)
	}
	removePerson()
}

func TestCreateIncorrectNamePerson(t *testing.T) {
	//post /users
	user := "{\"name\":23,\"age\":10}\n"
	response := postRequest(t, "/api-v1/user", strings.NewReader(user))
	checkResponseCode(t, http.StatusInternalServerError, response.Code)
	expected := ""
	if body := response.Body.String(); body != expected {
		t.Errorf("Expected empty but. Got %s", body)
	}
	removePerson()
}

func TestCreateIncorrectAgePerson(t *testing.T) {
	//post /users
	user := "{\"name\":\"Dani\",\"age\":\"a\"}\n"
	response := postRequest(t, "/api-v1/user", strings.NewReader(user))
	checkResponseCode(t, http.StatusInternalServerError, response.Code)
	expected := ""
	if body := response.Body.String(); body != expected {
		t.Errorf("Expected empty but. Got %s", body)
	}
	removePerson()
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
