package controller

import (
	"Project_BNCC_GO/model"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {
	user := model.User{}
	c.Bind(&user)

	// ganti pake jwt
	return c.String(http.StatusOK, "LOGIN")
}

func SignUP() {

}
