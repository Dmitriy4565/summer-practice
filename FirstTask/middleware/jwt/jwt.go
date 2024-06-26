package jwt

import (
    "net/http"
    "github.com/labstack/echo/v4"
    "github.com/golang-jwt/jwt"
	"FirstTask/models"
)

var jwtSecret = []byte("supersecretkey")

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        token := c.Request().Header.Get("Authorization")
        if token == "" {
            return echo.NewHTTPError(http.StatusUnauthorized, "Missing JWT token")
        }
        
        claims := &models.JWTCustomClaims{}
        t, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
            return jwtSecret, nil
        })
        if err != nil {
            return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
        }
        if !t.Valid {
            return echo.NewHTTPError(http.StatusUnauthorized, "Invalid JWT token")
        }

        c.Set("user", t)

        return next(c)
    }
}
