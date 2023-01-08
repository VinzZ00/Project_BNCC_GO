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

func AddMemory(c echo.Context) error {
	return c.Render(http.StatusOK, "add-memory", nil)
}

func MemoryDetail(c echo.Context) error {
	id := c.Param("id")
	return c.Render(http.StatusOK, "memory-detail", id)
}