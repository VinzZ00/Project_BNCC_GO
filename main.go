package main

import (
	"Project_BNCC_GO/config"
	"Project_BNCC_GO/controller"
	"Project_BNCC_GO/handler"
	"Project_BNCC_GO/utils"
	"errors"
	"html/template"
	"io"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
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
		SigningKey:  config.JWT_KEY,
		TokenLookup: "cookie:Token",
		ErrorHandler: func(c echo.Context, err error) error {
			return utils.SendResponse(c, utils.BaseResponse{
				StatusCode: http.StatusUnauthorized,
				Message:    "You are not authorized",
			})
		},
		ParseTokenFunc: func(c echo.Context, auth string) (interface{}, error) {
			token, err := jwt.ParseWithClaims(auth, &config.JwtClaim{}, func(t *jwt.Token) (interface{}, error) {
				return config.JWT_KEY, nil
			})

			if err != nil {
				return nil, err
			}
			if !token.Valid {
				return nil, errors.New("invalid token")
			}

			return token, nil
		},
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(config.JwtClaim)
		},
	}))

	memoryGroup.POST("", controller.CreateMemory)
	memoryGroup.PUT("/:id", controller.UpdateMemory)
	memoryGroup.DELETE("/:id", controller.DeleteMemory)
	memoryGroup.GET("", controller.GetAllMemories)
	memoryGroup.GET("/:id", controller.GetAMemories)
	memoryGroup.DELETE("/:id/:picture_id", controller.DeletePicture)

	if err := e.Start(":5566"); err != nil {
		panic(err)
	}
}
