package handler

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

type M map[string]interface{}

func Home(c echo.Context) error {
	data := M{"message": "Hello World!"}
	return c.Render(http.StatusOK, "index", data)
}

func Login(c echo.Context) error {
	return c.Render(http.StatusOK, "login", nil)
}

func Register(c echo.Context) error {
	return c.Render(http.StatusOK, "register", nil)
}

func Memory(c echo.Context) error {
	return c.Render(http.StatusOK, "memories", nil)
}