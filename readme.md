# CarRentelApi

CarRentelApi is a simple car rental API implemented in Go, using Gorilla mux for routing and GORM for database interactions.

## Prerequisites

Before running the application, make sure you have the following prerequisites installed:

- Go (version X.X.X): [Installation Guide](https://golang.org/doc/install)
- MySQL Database: [Download MySQL](https://dev.mysql.com/downloads/mysql/)

## Setup

1. **Clone the repository:**

    ```bash
    git clone https://github.com/harim12/CarRentelApi.git
    cd CarRentelApi
    ```

2. **Database Configuration:**

    Create a MySQL database and update the database configuration in `main.go`:

    ```go
    const DNS = "root:admin@tcp(127.0.0.1:3306)/carrenteldb?charset=utf8mb4&parseTime=True&loc=Local"
    ```

    Update the database credentials and connection details accordingly.

3. **Build and Run:**

    ```bash
    go build
    ./CarRentelApi.exe
    ```

4. **Access the API:**

    The API will be running at `http://localhost:9000`.

## API Endpoints

- **GET /cars:** Retrieve the list of cars.
- **POST /cars:** Create a new car.
- **POST /cars/{registration}/rentals:** Rent a car.
- **POST /cars/{registration}/returns:** Return a rented car.

