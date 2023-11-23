// car_test.go
package main

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"encoding/json"
	"os"
	"github.com/stretchr/testify/assert"

	"github.com/gorilla/mux"
)

func setupRouter() *mux.Router {
    return globalRouter
}
const TestDNS = "root:admin@tcp(127.0.0.1:3306)/carrentaltestdb?charset=utf8mb4&parseTime=True&loc=Local"


func TestMain(m *testing.M) {
    // Perform setup here, including initializing globalRouter

	initializeRouter()
    // Run tests
    result := m.Run()

    // Exit with the result of the tests
    os.Exit(result)
}


func resetDB(t *testing.T) {
	// Reset the database to its initial state
	db, err := sql.Open("mysql", TestDNS)
	if err != nil {
		t.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("DROP TABLE IF EXISTS cars")
	if err != nil {
		t.Fatalf("Error dropping cars table: %v", err)
	}

	initQuery := `
		CREATE TABLE IF NOT EXISTS cars (
			model VARCHAR(255),
			registration VARCHAR(255) PRIMARY KEY,
			mileage INT,
			available BOOLEAN
		)
	`

	_, err = db.Exec(initQuery)
	if err != nil {
		t.Fatalf("Error creating cars table: %v", err)
	}
}



func TestListCars(t *testing.T) {
	// Initialize the database and server for testing
	resetDB(t)
	r := setupRouter()

	// Create a new car
	payload := []byte(`{"model": "TestCar", "registration": "Testing1234", "mileage": 0}`)
	req, err := http.NewRequest("POST", "/cars", bytes.NewBuffer(payload))
	assert.Nil(t, err, "Error creating request")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Unexpected status code for car creation")

	req, err = http.NewRequest("GET", "/cars", nil)
	assert.Nil(t, err, "Error creating get cars request")

	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Unexpected status code for get cars request")

	// Decode the response body
	var cars []map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &cars)
	assert.Nil(t, err, "Error decoding response body")

	// Check if the list of cars contains the expected car
	found := false
	for _, car := range cars {
		if car["registration"] == "Testing1234" {
			found = true
			break
		}
	}
	assert.True(t, found, "Expected car not found in the list")
}

func TestCreateCar(t *testing.T) {
	// Initialize the database and server for testing
	resetDB(t)

	r := setupRouter()

	// Create a new car
	payload := []byte(`{"model": "TestCar", "registration": "TESTE1234", "mileage": 0}`)
	req, err := http.NewRequest("POST", "/cars", bytes.NewBuffer(payload))
	assert.Nil(t, err, "Error creating request")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code, "Unexpected status code")

	// Decode the response body
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Nil(t, err, "Error decoding response body")

	// Check the message and ID in the response
	assert.Equal(t, "Car added successfully", response["message"])
	assert.NotNil(t, response["id"])

}

func TestRentCar(t *testing.T) {
	// Initialize the database and server for testing
	resetDB(t)
	r := setupRouter()

	// Create a new car
	payload := []byte(`{"model": "TestCar", "registration": "TESTE1234", "mileage": 0}`)
	req, err := http.NewRequest("POST", "/cars", bytes.NewBuffer(payload))
	assert.Nil(t, err, "Error creating request")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Rent the created car
	req, err = http.NewRequest("POST", "/cars/TESTE1234/rentals", nil)
	assert.Nil(t, err, "Error creating rental request")

	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code, "Unexpected status code for renting a car")
}

func TestReturnCar(t *testing.T) {
	resetDB(t)
	r := setupRouter()

	// Create a new car
	payload := []byte(`{"model": "TestCar", "registration": "TESTE1234", "mileage": 0}`)
	req, err := http.NewRequest("POST", "/cars", bytes.NewBuffer(payload))
	assert.Nil(t, err, "Error creating request")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Rent the created car
	req, err = http.NewRequest("POST", "/cars/TESTE1234/rentals", nil)
	assert.Nil(t, err, "Error creating rental request")

	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Return the rented car
	returnPayload := []byte(`{"kilometers": 50}`)
	req, err = http.NewRequest("POST", "/cars/TESTE1234/returns", bytes.NewBuffer(returnPayload))
	assert.Nil(t, err, "Error creating return car request")

	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code, "Unexpected status code for returning a car")
}
