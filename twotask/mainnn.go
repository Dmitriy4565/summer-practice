package main

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JWTCustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var users = make(map[string]string)

var jwtSecret = []byte("supersecretkey")

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/login", login)
	e.POST("/register", register)

	r := e.Group("/restricted")
	r.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    jwtSecret,
		SigningMethod: "HS256",
	}))
	r.GET("/data", restrictedData)
	r.GET("/info", restrictedInfo)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

func register(c echo.Context) error {
	u := new(User)
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

func login(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Invalid input",
		})
	}
	password, ok := users[u.Username]
	if !ok || password != u.Password {
		return echo.ErrUnauthorized
	}

	claims := &JWTCustomClaims{
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

func restrictedData(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Welcome " + username + "!",
	})
}

func restrictedInfo(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	return c.JSON(http.StatusOK, echo.Map{
		"info": "This is secret information for " + username,
	})
}
