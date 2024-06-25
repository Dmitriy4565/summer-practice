package handlers

import (
    "net/http"
    "github.com/golang-jwt/jwt/v4"
    "github.com/labstack/echo/v4"
)

func RestrictedData(c echo.Context) error {
    user := c.Get("user").(*jwt.Token)
    claims := user.Claims.(jwt.MapClaims)
    username := claims["username"].(string)
    return c.JSON(http.StatusOK, echo.Map{
        "message": "Welcome " + username + "!",
    })
}

func RestrictedInfo(c echo.Context) error {
    user := c.Get("user").(*jwt.Token)
    claims := user.Claims.(jwt.MapClaims)
    username := claims["username"].(string)
    return c.JSON(http.StatusOK, echo.Map{
        "info": "This is secret information for " + username,
    })
}