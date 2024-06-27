package main

import (
    "database/sql"
    "net/http"

    "github.com/labstack/echo/v4"
    _ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
    var err error
    connStr := "user=postgres password=12345678 dbname=cars sslmode=disable"
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        panic(err)
    }
}

// Получить все машины пользователя
func getUserCars(c echo.Context) error {
    userID := c.Param("id")
    rows, err := db.Query("SELECT id, make, model FROM cars WHERE user_id = $1", userID)
    if err != nil {
        return err
    }
    defer rows.Close()

    var cars []map[string]interface{}
    for rows.Next() {
        var id int
        var make, model string
        if err := rows.Scan(&id, &make, &model); err != nil {
            return err
        }
        cars = append(cars, map[string]interface{}{"id": id, "make": make, "model": model})
    }

    return c.JSON(http.StatusOK, cars)
}

// Получить все двигатели пользователя
func getUserEngines(c echo.Context) error {
    userID := c.Param("id")
    rows, err := db.Query(`
        SELECT DISTINCT e.id, e.type, e.horsepower
        FROM engines e
        JOIN cars c ON e.car_id = c.id
        WHERE c.user_id = $1
    `, userID)
    if err != nil {
        return err
    }
    defer rows.Close()

    var engines []map[string]interface{}
    for rows.Next() {
        var id int
        var typ string
        var horsepower int
        if err := rows.Scan(&id, &typ, &horsepower); err != nil {
            return err
        }
        engines = append(engines, map[string]interface{}{"id": id, "type": typ, "horsepower": horsepower})
    }

    return c.JSON(http.StatusOK, engines)
}

// Получить двигатель конкретной машины
func getCarEngine(c echo.Context) error {
    carID := c.Param("id")
    row := db.QueryRow("SELECT id, type, horsepower FROM engines WHERE car_id = $1", carID)

    var id int
    var typ string
    var horsepower int
    if err := row.Scan(&id, &typ, &horsepower); err != nil {
        if err == sql.ErrNoRows {
            return c.JSON(http.StatusNotFound, map[string]string{"message": "Engine not found"})
        }
        return err
    }

    return c.JSON(http.StatusOK, map[string]interface{}{"id": id, "type": typ, "horsepower": horsepower})
}

// Получить все двигатели по марке машины
func getEngines(c echo.Context) error {
    make := c.Param("make")
    rows, err := db.Query(`
        SELECT e.id, e.type, e.horsepower
        FROM engines e
        JOIN cars c ON e.car_id = c.id
        WHERE c.make = $1
    `, make)
    if err != nil {
        return err
    }
    defer rows.Close()

    var engines []map[string]interface{}
    for rows.Next() {
        var id int
        var typ string
        var horsepower int
        if err := rows.Scan(&id, &typ, &horsepower); err != nil {
            return err
        }
        engines = append(engines, map[string]interface{}{"id": id, "type": typ, "horsepower": horsepower})
    }

    return c.JSON(http.StatusOK, engines)
}

func main() {
    initDB()
    defer db.Close()

    e := echo.New()

    e.GET("/user/:id/cars", getUserCars)
    e.GET("/user/:id/engines", getUserEngines)
    e.GET("/car/:id/engine", getCarEngine)
    e.GET("/engines/:make", getEngines)

    e.Start(":8080")
}