package main

import (
	"errors"
)

// ListCarsFromDB fetches and returns the list of cars from the database
func ListCarsFromDB() ([]Car, error) {
	rows, err := db.Query("SELECT model, registration, mileage, available FROM cars")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cars []Car
	for rows.Next() {
		var car Car
		err := rows.Scan(&car.Model, &car.Registration, &car.Mileage, &car.Available)
		if err != nil {
			return nil, err
		}
		cars = append(cars, car)
	}

	return cars, nil
}

// AddCarToDB adds a new car to the database
func AddCarToDB(newCar Car) (int64, error) {
	result, err := db.Exec("INSERT INTO cars (model, registration, mileage, available) VALUES (?, ?, ?, ?)",
		newCar.Model, newCar.Registration, newCar.Mileage, true)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// RentCarInDB marks a car as rented in the database
func RentCarInDB(registration string) error {
	result, err := db.Exec("UPDATE cars SET available = false WHERE registration = ?", registration)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("car not found")
	}

	return nil
}

// ReturnCarInDB marks a rented car as returned in the database and updates the mileage
func ReturnCarInDB(registration string, kilometers int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Check if the car is rented
	var available bool
	err = tx.QueryRow("SELECT available FROM cars WHERE registration = ?", registration).Scan(&available)
	if err != nil {
		tx.Rollback()
		return err
	}

	if !available {
		// Update mileage and mark as available
		_, err := tx.Exec("UPDATE cars SET mileage = mileage + ?, available = true WHERE registration = ?", kilometers, registration)
		if err != nil {
			tx.Rollback()
			return err
		}

		err = tx.Commit()
		if err != nil {
			return err
		}

		return nil
	}

	tx.Rollback()
	return errors.New("car not rented")
}
