package main

import (
	"Project_BNCC_GO/controller"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		fmt.Println("Hello")
		return c.String(http.StatusOK, "Hello Welcome to this PAGE")
	})

	e.POST("/login", controller.Login)

	e.POST("/createMemo", controller.CreateMemory)

	if err := e.Start(":5566"); err != nil {
		panic(err)
	}

}
