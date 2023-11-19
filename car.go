package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error


// DNS is the Data Source Name for connecting to the MySQL database

const DNS = "root:admin@tcp(127.0.0.1:3306)/carrenteldb?charset=utf8mb4&parseTime=True&loc=Local"

// Car is the model for representing a car in the database
type Car struct {
	gorm.Model
	ModeleCar        string `json:"model"`
	Registration string `json:"registration"`
	Mileage      int    `json:"mileage"`
	Available    bool   `json:"available" gorm:"default:true"`
}

// InitialMigration performs the initial migration of database schema
func InitialMigration() {
	DB, err = gorm.Open(mysql.Open(DNS), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Cannot connect to DB")
	}
	DB.AutoMigrate(&Car{})
}

// GetCars handles the HTTP GET request to retrieve a list of cars
func GetCars(w http.ResponseWriter, r *http.Request) {
	var cars []Car
	DB.Find(&cars)
	response := make([]map[string]interface{}, len(cars))
	for i, car := range cars {
		response[i] = map[string]interface{}{
			"model":        car.ModeleCar,
			"registration": car.Registration,
			"mileage":      car.Mileage,
			"available":    car.Available,
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateCar handles the HTTP POST request to create a new car
func CreateCar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var car Car
	json.NewDecoder(r.Body).Decode(&car)

	// Check if the car already exists
	var existingCar Car
	result := DB.Where("registration = ?", car.Registration).First(&existingCar)
	if result.Error == nil {
		http.Error(w, "Car with the same registration number already exists", http.StatusBadRequest)
		return
	}

	DB.Create(&car)
	json.NewEncoder(w).Encode(car)
}

// RentCar handles the HTTP POST request to mark a car as rented
func RentCar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	registration := params["registration"]

	var car Car
	result := DB.Where("registration = ?", registration).First(&car)

	if result.Error != nil {
		http.Error(w, "Car not found", http.StatusNotFound)
		return
	}

	if !car.Available {
		http.Error(w, "Car is already rented", http.StatusBadRequest)
		return
	}

	car.Available = false
	DB.Save(&car)
	w.WriteHeader(http.StatusOK)
}

// ReturnCar handles the HTTP POST request to mark a car as returned
func ReturnCar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	registration := params["registration"]

	var car Car
	result := DB.Where("registration = ?", registration).First(&car)

	if result.Error != nil {
		http.Error(w, "Car not found", http.StatusNotFound)
		return
	}

	if car.Available {
		http.Error(w, "Car is not marked as rented", http.StatusBadRequest)
		return
	}

	var returnInfo struct {
		Kilometers int `json:"kilometers"`
	}

	err := json.NewDecoder(r.Body).Decode(&returnInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	car.Mileage += returnInfo.Kilometers
	car.Available = true
	DB.Save(&car)
	w.WriteHeader(http.StatusOK)
}
