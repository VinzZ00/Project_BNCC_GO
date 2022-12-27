package controller

import (
	"Project_BNCC_GO/model"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {
	LoginWebData := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	if err := c.Bind(&LoginWebData); err != nil {
		panic(err)
	}

	user := model.User{
		Email:    LoginWebData.Email,
		Password: LoginWebData.Password,
	}

	result := db.Find(&user)

	if err := result.Error; err != nil {
		panic(err)
	}

	// JWT

	return c.String(http.StatusOK, "LOGIN")
}

func SignUP(c echo.Context) error {
	signupUser := model.User{}

	c.Bind(&signupUser)
	result := db.Create(&signupUser)

	if result.Error != nil {
		panic(result.Error)
	}

	return c.JSON(http.StatusOK, struct {
		message string
	}{
		message: "User with ID " + string(signupUser.Userid) + "is Created",
	})
}

func init() {
	if db == nil {
		if database, err := model.GetDB(); err == nil {
			db = database
		} else {
			panic(err)
		}
	}
}
