package main

import (
	"Project_BNCC_GO/config"
	"Project_BNCC_GO/controller"
	"fmt"
	"net/http"

	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		fmt.Println("Hello")
		return c.String(301, "Hello Welcome to this PAGE")
	})

	authGroup := e.Group("/auth")
	authGroup.POST("/login", controller.Login)
	authGroup.POST("/register", controller.SignUP)

	memoryGroup := e.Group("/memories")
	memoryGroup.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWT_KEY),
		TokenLookup: "cookie:Token",
		ErrorHandler: func(c echo.Context, err error) error {
			c.JSON(http.StatusUnauthorized, struct {
				Message string
			}{
				Message: "Status Un Authorized",
			})
			return nil
		},
	}))
	memoryGroup.GET("/test", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct {
			Message string
		}{
			Message: "hellowFromProtectedAPI, Ini Harus dalam kondisi sudah signed In",
		})
	})

	memoryGroup.POST("/", controller.CreateMemory)


	if err := e.Start(":5566"); err != nil {
		panic(err)
	}
}
