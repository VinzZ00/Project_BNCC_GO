package controller

import (
	"Project_BNCC_GO/model"
	"Project_BNCC_GO/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func DeletePicture(c echo.Context) error {
	params := struct {
		MemoryIDParam
		pictureId string `param:"picture_id" validate:"required,number"`
	}{}

	c.Bind(&params)
	deletedPic := model.Picture{}
	if err := db.Unscoped().Where("id = ? And memory_id = ?", params.pictureId, params.ID).Delete(&deletedPic).Error; err != nil {
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