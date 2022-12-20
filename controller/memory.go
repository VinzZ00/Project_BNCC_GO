package controller

import (
	"Project_BNCC_GO/model"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var db *gorm.DB

func CreateMemory(c echo.Context) error {
	m := new(model.Memory)

	if err := c.Bind(&m); err != nil {
		panic(err)
	}
	m.DateModified = time.Now()
	fmt.Println(m)
	result := db.Create(m)

	fmt.Println(result)

	return c.Redirect(http.StatusOK, "localhost:5555")
}

func init() {
	if database, err := model.GetDB(); err == nil {
		db = database
	} else {
		panic(err)
	}

}
