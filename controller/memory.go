package controller

import (
	"Project_BNCC_GO/model"
	"Project_BNCC_GO/utils"
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/png"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var db *gorm.DB

type MemoryIDParam struct {
	ID int `param:"id" validate:"required,number"`
}

func CreateMemory(c echo.Context) error {
	// Struct untuk ambil data web
	body := struct {
		Description string            `json:"description"`
		UserId      uint              `json:"userId"`
		Paths       []string          `json:"picturePaths"`
		Tags        []model.MemoryTag `json:"tags"`
	}{}
	if err := c.Bind(&body); err != nil {
		panic("Error di binding data")
	}

	picturesBytes := [][]byte{}
	for _, value := range body.Paths {
		imageFile, err := os.Open(value)
		if err != nil {
			panic("path invalid")
		}

		imageData, _, err := image.Decode(imageFile)
		if err != nil {
			panic(err)
		}

		buff := new(bytes.Buffer)
		if err = png.Encode(buff, imageData); err != nil {
			panic(err)
		}
		picturesBytes = append(picturesBytes, buff.Bytes())
	}

	pictures := []model.Picture{}
	// todo: karna modelnya diganti biar pake `string` untuk nerima base64,
	//       jadinya code dibawah gk bisa dipake
	//
	// for _, value := range picturesBytes {
	// 	pic := model.Picture{
	// 		Data: value,
	// 	}
	// 	pictures = append(pictures, pic)
	// }

	memory := model.Memory{
		Description:  body.Description,
		UserID:       body.UserId,
		MemoriesTags: body.Tags,
		Pictures:     pictures,
	}
	if err := db.Create(&memory).Error; err != nil {
		panic(err)
	}

	return utils.SendResponse(c, utils.BaseResponse{
		StatusCode: http.StatusCreated,
		Message:    "Memory is successfully created",
	})
}

func UpdateMemory(c echo.Context) error {
	payload := struct {
		MemoryIDParam
		Description string      `json:"description"`
		Tags        []model.Tag `json:"tags"`
	}{}

	if err := c.Bind(&payload); err != nil {
		return utils.SendResponse(c, utils.BaseResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
	}

	memoryId := payload.ID
	var memory model.Memory
	if err := db.First(&memory, memoryId).Error; err != nil {
		return utils.SendResponse(c, utils.BaseResponse{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
		})
	}

	currentUser, _ := utils.GetAuthUser(c)
	if memory.UserID != currentUser.UserID {
		return utils.SendResponse(c, utils.BaseResponse{
			StatusCode: http.StatusForbidden,
			Message:    "You do not own this memory",
		})
	}

	// membuat tag yang terupdate baru dari response web
	tags := []model.MemoryTag{}
	for _, tag := range payload.Tags {
		tags = append(tags, model.MemoryTag{
			MemoryID: memory.ID,
			TagID:    tag.ID,
		})
	}

	memory.Description = payload.Description
	memory.UpdatedAt = time.Now()
	memory.MemoriesTags = tags

	if err := db.Save(&memory).Error; err != nil {
		return utils.SendResponse(c, utils.BaseResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	return utils.SendResponse(c, utils.BaseResponse{
		StatusCode: http.StatusOK,
		Message:    "Memory has been updated",
	})
}

func DeleteMemory(c echo.Context) error {
	payload := MemoryIDParam{}
	if err := c.Bind(&payload); err != nil {
		return utils.SendResponse(c, utils.BaseResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
	}

	memoryId := payload.ID
	var memory model.Memory
	if err := db.First(&memory, memoryId).Error; err != nil {
		response := utils.BaseResponse{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.StatusCode = http.StatusNotFound
		}

		return utils.SendResponse(c, response)
	}

	currentUser, _ := utils.GetAuthUser(c)
	if memory.UserID != currentUser.UserID {
		return utils.SendResponse(c, utils.BaseResponse{
			StatusCode: http.StatusForbidden,
			Message:    "You do not own this memory",
		})
	}
	if err := db.Delete(&memory).Error; err != nil {
		return utils.SendResponse(c, utils.BaseResponse{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		})
	}

	return utils.SendResponse(c, utils.BaseResponse{
		StatusCode: http.StatusAccepted,
		Message:    fmt.Sprintf("Memory with ID %d has been deleted", memoryId),
	})
}

func init() {
	if database, err := model.GetDB(); err == nil {
		db = database
	} else {
		panic(err)
	}
}
