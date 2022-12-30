package main

import (
	"Project_BNCC_GO/config"
	"Project_BNCC_GO/controller"
	"fmt"
	"net/http"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.RemoveTrailingSlash())
	e.GET("/", func(c echo.Context) error {
		fmt.Println("Hello")
		return c.String(301, "Hello Welcome to this PAGE")
	})

	authGroup := e.Group("/auth")
	authGroup.POST("/login", controller.Login)
	authGroup.POST("/register", controller.SignUP)

	memoryGroup := e.Group("/memories")
	memoryGroup.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(config.JWT_KEY),
		TokenLookup: "cookie:Token",
		ErrorHandler: func(c echo.Context, err error) error {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "You are not authorized",
			})
		},
	}))
	memoryGroup.GET("/test", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "hellowFromProtectedAPI, Ini Harus dalam kondisi sudah signed In",
		})
	})

	memoryGroup.POST("", controller.CreateMemory)
	memoryGroup.PUT("/:id", controller.UpdateMemory)

	if err := e.Start(":5566"); err != nil {
		panic(err)
	}
}
