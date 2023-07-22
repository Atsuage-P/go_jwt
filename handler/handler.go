package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func SignUp(c echo.Context) error {
	return c.String(http.StatusOK, "Signup")
}

func Login(c echo.Context) error {
	return c.String(http.StatusOK, "Login")
}

func Logout(c echo.Context) error {
	return c.String(http.StatusOK, "Logout")
}