package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/go-sql-driver/mysql"
)

// Car represents the car entity
type Car struct {
	Model        string `json:"model"`
	Registration string `json:"registration"`
	Mileage      int    `json:"mileage"`
	Available    bool   `json:"available"`
}

var db *sql.DB
var globalRouter *mux.Router

const DNS = "root:admin@tcp(127.0.0.1:3306)/carrenteldb?charset=utf8mb4&parseTime=True&loc=Local"

func init() {
	var err error
	db, err = sql.Open("mysql",DNS)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Create cars table if not exists
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS cars (
			model VARCHAR(255),
			registration VARCHAR(255) PRIMARY KEY,
			mileage INT,
			available BOOLEAN
		)
	`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}

func initializeRouter() {
	// Create a new Gorilla mux router
	r := mux.NewRouter()

	// Define routes and their corresponding handler functions
	r.HandleFunc("/cars", ListCars).Methods("GET")
	r.HandleFunc("/cars", AddCar).Methods("POST")
	r.HandleFunc("/cars/{registration}/rentals", RentCar).Methods("POST")
	r.HandleFunc("/cars/{registration}/returns", ReturnCar).Methods("POST")

	// Start the HTTP server on port 9000

	globalRouter = r
	
}


func main() {
	
	initializeRouter()
	log.Fatal(http.ListenAndServe(":9000", globalRouter))

}

// ListCars returns a list of all cars
func ListCars(w http.ResponseWriter, r *http.Request) {
	cars, err := ListCarsFromDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cars)
}

// AddCar adds a new car
func AddCar(w http.ResponseWriter, r *http.Request) {
	var newCar Car
	err := json.NewDecoder(r.Body).Decode(&newCar)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := AddCarToDB(newCar)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{"message": "Car added successfully", "id": id}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RentCar rents a car
func RentCar(w http.ResponseWriter, r *http.Request) {
	registration := mux.Vars(r)["registration"]

	err := RentCarInDB(registration)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Car rented successfully"}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ReturnCar returns a rented car
func ReturnCar(w http.ResponseWriter, r *http.Request) {
	registration := mux.Vars(r)["registration"]

	var returnDetails map[string]int
	err := json.NewDecoder(r.Body).Decode(&returnDetails)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ReturnCarInDB(registration, returnDetails["kilometers"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Car returned successfully"}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
