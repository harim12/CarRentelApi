# Car Rental API

This is a simple Car Rental API built with Go and GORM.

## Prerequisites

Before you begin, ensure you have the following:

- Go installed on your machine.
- MySQL database available.

## Setup

1. **Clone the repository:**

    ```bash
    git clone https://github.com/your-username/CarRentelApi.git
    ```

2. **Navigate to the project directory:**

    ```bash
    cd CarRentelApi
    ```

3. **Initialize the MySQL database:**

   - Create a MySQL database with the name `carrenteldb`.
   - Update the database connection details in `car.go`:

     ```go
     const DNS = "root:password@tcp(127.0.0.1:3306)/carrenteldb?charset=utf8mb4&parseTime=True&loc=Local"
     ```

     Replace `password` with your MySQL password.

4. **Run the application:**

    ```bash
    go run main.go
    ```

   The API server should now be running on [http://localhost:9000](http://localhost:9000).

## API Endpoints

- `GET /cars`: Get a list of cars.
- `POST /cars`: Create a new car.
- `POST /cars/{registration}/rentals`: Rent a car.
- `POST /cars/{registration}/returns`: Return a rented car.

Feel free to explore and use the provided API endpoints!
