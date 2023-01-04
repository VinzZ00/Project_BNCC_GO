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