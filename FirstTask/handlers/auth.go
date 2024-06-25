package handlers

import (
    "net/http"
    "time"
    "FirstTask/models"
    "github.com/golang-jwt/jwt/v4"
    "github.com/labstack/echo/v4"
)

var users = make(map[string]string)
var jwtSecret = []byte("supersecretkey")

func Register(c echo.Context) error {
    u := new(models.User)
    if err := c.Bind(u); err != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{
            "message": "Invalid input",
        })
    }
    if _, ok := users[u.Username]; ok {
        return c.JSON(http.StatusConflict, echo.Map{
            "message": "User already exists",
        })
    }
    users[u.Username] = u.Password
    return c.JSON(http.StatusOK, echo.Map{
        "message": "User registered",
    })
}

func Login(c echo.Context) error {
    u := new(models.User)
    if err := c.Bind(u); err != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{
            "message": "Invalid input",
        })
    }
    password, ok := users[u.Username]
    if !ok || password != u.Password {
        return echo.ErrUnauthorized
    }

    claims := &models.JWTCustomClaims{
        Username: u.Username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    t, err := token.SignedString(jwtSecret)
    if err != nil {
        return err
    }

    return c.JSON(http.StatusOK, echo.Map{
        "token": t,
    })
}