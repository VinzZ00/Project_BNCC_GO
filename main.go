package main

import (
	"Project_BNCC_GO/controller"
	"Project_BNCC_GO/handler"
	"Project_BNCC_GO/utils"
	"html/template"
	"io"

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

	authMiddleware := echojwt.WithConfig(utils.GetEchoJwtConfig())

	authGroup := e.Group("/auth")
	authGroup.POST("/login", controller.Login)
	authGroup.POST("/register", controller.SignUP)

	memoryGroup := e.Group("/memories")
	memoryGroup.Use(authMiddleware)

	memoryGroup.POST("", controller.CreateMemory)
	memoryGroup.PUT("/:id", controller.UpdateMemory)
	memoryGroup.DELETE("/:id", controller.DeleteMemory)
	memoryGroup.GET("", controller.GetAllMemories)
	memoryGroup.GET("/:id", controller.GetAMemories)

	memoryGroup.GET("/sort", controller.GetMemorySortBy)
	// memoryGroup.GET("sort", controller.GetMemorySortBy)
	// memoryGroup.GET("sort", controller.GetMemorySortBy)
	// memoryGroup.GET("sort", controller.GetMemorySortBy)

	pictureGroup := e.Group("/pictures")
	pictureGroup.DELETE("/:id", controller.DeletePicture, authMiddleware)

	if err := e.Start(":5566"); err != nil {
		panic(err)
	}
}
