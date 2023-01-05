package controller

import (
	"Project_BNCC_GO/model"
	"Project_BNCC_GO/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type PictureIDParam struct {
	ID int `param:"id" validate:"required,number"`
}

func DeletePicture(c echo.Context) error {
	params := PictureIDParam{}

	c.Bind(&params)
	if err := db.Delete(&model.Picture{}, params.ID).Error; err != nil {
		return utils.SendResponse(c, utils.BaseResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
	}

	return utils.SendResponse(c, utils.BaseResponse{
		StatusCode: http.StatusOK,
		Message:    "Picture Successfully deleted",
	})

}

func init() {
	if database, err := model.GetDB(); err == nil {
		db = database
	} else {
		panic(err)
	}
}
