package main

import (
	"Project_BNCC_GO/config"
	"Project_BNCC_GO/controller"
	"Project_BNCC_GO/handler"
	"Project_BNCC_GO/utils"
	"html/template"
	"io"

	"github.com/go-playground/validator/v10"
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
	e.Validator = &config.DefaultValidator{
		Validator: validator.New(),
	}
	e.Static("/images", "public/images")
	e.Static("/css", "public/css")
	e.Static("/js", "public/js")

	e.Use(middleware.RemoveTrailingSlash())
	e.GET("/", handler.Home)
	e.GET("/login", handler.Login)
	e.GET("/register", handler.Register)
	e.GET("/memories", handler.Memory)
	e.GET("/memories/add", handler.AddMemory)
	e.GET("/memories/:id", handler.MemoryDetail)

	authMiddleware := echojwt.WithConfig(utils.GetEchoJwtConfig())

	authGroup := e.Group("/api/auth")
	authGroup.POST("/login", controller.Login)
	authGroup.POST("/register", controller.SignUP)
	authGroup.DELETE("/logout", controller.Logout)

	memoryGroup := e.Group("/api/memories")
	memoryGroup.Use(authMiddleware)

	memoryGroup.POST("", controller.CreateMemory)
	memoryGroup.PUT("/:id", controller.UpdateMemory)
	memoryGroup.DELETE("/:id", controller.DeleteMemory)
	memoryGroup.GET("", controller.GetAllMemories)
	memoryGroup.GET("/:id", controller.GetAMemories)

	memoryGroup.GET("", controller.GetAllMemories)
	memoryGroup.GET("/filter", controller.MemoryFilterBy)
	memoryGroup.GET("/sort", controller.GetMemorySortBy)
	memoryGroup.POST("/:id", controller.AddPicture)
	memoryGroup.GET("/tags", controller.GetAllTags)
	// memoryGroup.GET("sort", controller.GetMemorySortBy)
	// memoryGroup.GET("sort", controller.GetMemorySortBy)

	pictureGroup := e.Group("/api/pictures")
	pictureGroup.Use(authMiddleware)
	pictureGroup.GET("/:id", controller.ReadPicture)
	pictureGroup.DELETE("/:id", controller.DeletePicture, authMiddleware)

	if err := e.Start(":5566"); err != nil {
		panic(err)
	}
}
