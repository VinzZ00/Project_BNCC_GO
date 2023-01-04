package main

import (
	"Project_BNCC_GO/config"
	"Project_BNCC_GO/controller"
	"Project_BNCC_GO/handler"
	"html/template"
	"io"
	"Project_BNCC_GO/utils"
	"net/http"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
    templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()
	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
	e.Renderer = t
	e.Static("/images", "public/images")
	
	e.Use(middleware.RemoveTrailingSlash())
	e.GET("/", handler.Home)
	e.GET("/login", handler.Login)
	e.GET("/register", handler.Register)

	authGroup := e.Group("/auth")
	authGroup.POST("/login", controller.Login)
	authGroup.POST("/register", controller.SignUP)

	memoryGroup := e.Group("/memories")
	memoryGroup.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(config.JWT_KEY),
		TokenLookup: "cookie:Token",
		ErrorHandler: func(c echo.Context, err error) error {
			return utils.SendResponse(c, utils.BaseResponse{
				StatusCode: http.StatusUnauthorized,
				Message: "You are not authorized",
			})
		},
	}))

	memoryGroup.POST("", controller.CreateMemory)
	memoryGroup.PUT("/:id", controller.UpdateMemory)
	memoryGroup.DELETE("/:id", controller.DeleteMemory)

	if err := e.Start(":5566"); err != nil {
		panic(err)
	}
}
