// car_test.go
package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"

	"github.com/gorilla/mux"
)

func setupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/cars", GetCars).Methods("GET")
	r.HandleFunc("/cars", CreateCar).Methods("POST")
	r.HandleFunc("/cars/{registration}/rentals", RentCar).Methods("POST")
	r.HandleFunc("/cars/{registration}/returns", ReturnCar).Methods("POST")
	return r
}


func TestCreateCar(t *testing.T) {
	// Initialize the database and server for testing
	InitialMigration()
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
	var car Car
	err = json.Unmarshal(rr.Body.Bytes(), &car)
	assert.Nil(t, err, "Error decoding response body")

	// Check if the car was created with the expected values
	assert.Equal(t, "TestCar", car.ModeleCar, "Unexpected car model")
	// Add more checks for other fields as needed
}

func TestGetCars(t *testing.T) {
	// Initialize the database and server for testing
	InitialMigration()
	r := setupRouter()

	// Create a new car
	payload := []byte(`{"model": "TestCar", "registration": "Testing1234", "mileage": 0}`)
	req, err := http.NewRequest("POST", "/cars", bytes.NewBuffer(payload))
	assert.Nil(t, err, "Error creating request")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Check the response status code for car creation
	assert.Equal(t, http.StatusOK, rr.Code, "Unexpected status code for car creation")

	// Get the list of cars
	req, err = http.NewRequest("GET", "/cars", nil)
	assert.Nil(t, err, "Error creating get cars request")

	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Check the response status code for get cars request
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

func TestRentCar(t *testing.T) {
	// Initialize the database and server for testing
	InitialMigration()
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
	assert.Equal(t, http.StatusOK, rr.Code, "Unexpected status code")
	// Add more checks for the rent operation as needed
}

func TestReturnCar(t *testing.T) {
	// Initialize the database and server for testing
	InitialMigration()
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
	assert.Equal(t, http.StatusOK, rr.Code, "Unexpected status code")
}