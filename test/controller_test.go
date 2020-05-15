package controllertest

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/rfaguiar/golang-first-api/api"
	"github.com/rfaguiar/golang-first-api/api/model"
	"github.com/rfaguiar/golang-first-api/database"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var server api.Server
var container testcontainers.Container

func TestMain(m *testing.M) {
	server = api.Server{}
	server.Initialize()

	setup()
	code := m.Run()

	tearDown()
	os.Exit(code)
}

func setup() {
	password := "Postgres2018!"
	username := "postgres"
	dbname := "service"
	var env = map[string]string{
		"POSTGRES_PASSWORD": password,
		"POSTGRES_USER":     username,
		"POSTGRES_DB":       dbname,
	}

	var portC = "5432/tcp"

	dbURL := func(port nat.Port) string {
		return fmt.Sprintf("postgres://%s:%s@localhost:%s/service?sslmode=disable",
			username, password, port.Port())
	}

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:latest",
			ExposedPorts: []string{portC},
			Cmd:          []string{"postgres", "-c", "fsync=off"},
			Env:          env,
			WaitingFor:   wait.ForSQL(nat.Port(portC), "postgres", dbURL).Timeout(time.Second * 5),
		},
		Started: true,
	}
	var err error
	container, err = testcontainers.GenericContainer(context.Background(), req)

	if err != nil {
		log.Fatal("start ", err)
	}
	mappedHost, err := container.Host(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	mappedPort, err := container.MappedPort(context.Background(), nat.Port(portC))
	if err != nil {
		log.Fatal(err)
	}

	os.Setenv("ENV_DB_HOST", mappedHost)
	os.Setenv("ENV_DB_PORT", mappedPort.Port())
	os.Setenv("ENV_DB_USER", username)
	os.Setenv("ENV_DB_PASSWORD", password)
	os.Setenv("ENV_DB_DBNAME", dbname)

	db := database.Open()
	database.ExecuteMigration(db)
}

func tearDown() {
	database.Close()
	container.Terminate(context.Background())
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	server.Router.ServeHTTP(rr, req)
	return rr
}
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestHealthCheck(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	server.Router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	var status model.Health
	if err = json.Unmarshal([]byte(rr.Body.String()), &status); err != nil {
		t.Errorf("Cannot convert to json: %v", err)
	}
	expected := model.Health{Status: "UP"}
	if status != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestEmptyTablePerson(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api-v1/user", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	if body := response.Body.String(); body != "[]\n" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}
