package controller

import (
	"Project_BNCC_GO/model"
	"Project_BNCC_GO/utils"
	"encoding/base64"
	"net/http"

	"github.com/labstack/echo/v4"
)

type PictureIDParam struct {
	ID int `param:"id" validate:"required,number"`
}

func ReadPicture(c echo.Context) error {
	params := PictureIDParam{}
	if err := c.Bind(&params); err != nil {
		return utils.SendResponse(c, utils.BaseResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
	}

	var picture model.Picture
	if err := db.First(&picture, params.ID).Error; err != nil {
		return utils.SendResponse(c, utils.BaseResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	decodedImage, err := base64.StdEncoding.DecodeString(picture.Data)
	if err != nil {
		return utils.SendResponse(c, utils.BaseResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
	}

	contentType := http.DetectContentType(decodedImage)
	return c.Blob(http.StatusOK, contentType, decodedImage)
}

func DeletePicture(c echo.Context) error {
	params := PictureIDParam{}

	if err := c.Bind(&params); err != nil {
		return utils.SendResponse(c, utils.BaseResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
	}

	var picture model.Picture
	var memory model.Memory
	currentUser, _ := utils.GetAuthUser(c)

	if err := db.First(&picture, params.ID).Error; err != nil {
		return utils.SendResponse(c, utils.BaseResponse{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
		})
	}

	// ignoring error because if the previous query succeeds
	// it means that it's parent (the memory) exist and the db is working fine.
	db.First(&memory, picture.MemoryID)

	if memory.UserID != currentUser.UserID {
		return utils.SendResponse(c, utils.BaseResponse{
			StatusCode: http.StatusForbidden,
			Message:    "You do not own this memory",
		})
	}

	if err := db.Delete(&model.Picture{}, params.ID).Error; err != nil {
		return utils.SendResponse(c, utils.BaseResponse{
			StatusCode: http.StatusInternalServerError,
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
