package controller_test

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/rfaguiar/golang-first-api/api"
	"github.com/rfaguiar/golang-first-api/database"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"io"
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

func getRequest(t *testing.T, url string) *httptest.ResponseRecorder {
	return executeRequest(t, http.MethodGet, url, nil)
}

func postRequest(t *testing.T, url string, body io.Reader) *httptest.ResponseRecorder {
	return executeRequest(t, http.MethodPost, url, body)
}

func deleteRequest(t *testing.T, url string) *httptest.ResponseRecorder {
	return executeRequest(t, http.MethodDelete, url, nil)
}

func putRequest(t *testing.T, url string, body io.Reader) *httptest.ResponseRecorder {
	return executeRequest(t, http.MethodPut, url, body)
}

func executeRequest(t *testing.T, method, url string, body io.Reader) *httptest.ResponseRecorder {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Errorf("Error: %s\n when execute %s %s", method, url, err)
	}
	rr := httptest.NewRecorder()
	server.Router.ServeHTTP(rr, request)
	return rr
}
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d but. Got %d\n", expected, actual)
	}
}
