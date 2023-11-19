package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func initializeRouter() {
	// Create a new Gorilla mux router
	r := mux.NewRouter()

	// Define routes and their corresponding handler functions
	r.HandleFunc("/cars", GetCars).Methods("GET")
	r.HandleFunc("/cars", CreateCar).Methods("POST")
	r.HandleFunc("/cars/{registration}/rentals", RentCar).Methods("POST")
	r.HandleFunc("/cars/{registration}/returns", ReturnCar).Methods("POST")

	// Start the HTTP server on port 9000
	log.Fatal(http.ListenAndServe(":9000", r))
}

func main() {
	// Perform the initial database migration
	InitialMigration()

	// Initialize the router
	initializeRouter()

	// Print a message to indicate that the server is running
	log.Println("Server is running on port 9000")
	
	// Start the HTTP server on port 9000
	log.Fatal(http.ListenAndServe(":9000", nil))
}
