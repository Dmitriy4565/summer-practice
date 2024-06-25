package main

import (
	"FirstTask/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/login", handlers.Login)
	e.POST("/register", handlers.Register)

	r := e.Group("/restricted")
	r.Use()
	r.GET("/data", handlers.RestrictedData)
	r.GET("/info", handlers.RestrictedInfo)

	e.Logger.Fatal(e.Start(":1323"))
}
