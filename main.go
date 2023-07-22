package main

import (
	"go_oauth/handler"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", handler.SignUp)
	e.Logger.Fatal(e.Start(":8080"))
}
