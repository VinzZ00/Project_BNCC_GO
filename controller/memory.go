package controller

import (
	"Project_BNCC_GO/model"
	"Project_BNCC_GO/utils"
	"errors"
	"fmt"
	"net/http"
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
	payload := struct {
		Description   string   `json:"description"`
		UserId        uint     `json:"userId"`
		base64Picture []string `json:"pictures"`
		Tags          []string `json:"tags"`
	}{}
	if err := c.Bind(&payload); err != nil {
		return utils.SendResponse(c, utils.BaseResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	// Pictures
	pictures := []model.Picture{}
	for _, val := range payload.base64Picture {
		picture := model.Picture{
			BaseModel: model.BaseModel{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Data: val,
		}
		pictures = append(pictures, picture)
	}

	//Tags
	tags := []model.Tag{}
	for _, val := range payload.Tags {
		checkTag := model.Tag{}

		err := db.Where("name= ?", val).First(&checkTag)
		if err != nil {
			tag := model.Tag{
				BaseModel: model.BaseModel{
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Name: val,
			}

			if err := db.Create(&tag).Error; err != nil {
				panic(err)
			}

			t := model.Tag{}
			db.Where("name = ? ", val).First(&tag)
			tags = append(tags, t)
		} else {
			tags = append(tags, checkTag)
		}
	}

	//memoryTag
	memoryTags := []model.MemoryTag{}

	for _, val := range tags {
		memoryTag := model.MemoryTag{
			TagID: val.ID,
		}
		memoryTags = append(memoryTags, memoryTag)
	}

	memory := model.Memory{
		BaseModel: model.BaseModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Description:  payload.Description,
		UserID:       payload.UserId,
		Pictures:     pictures,
		MemoriesTags: memoryTags,
	}

	if err := db.Create(&memory).Error; err != nil {
		panic(err)
	}

	// picturesBytes := [][]byte{}
	// for _, value := range payload.Paths {
	// 	imageFile, err := os.Open(value)
	// 	if err != nil {
	// 		panic("path invalid")
	// 	}

	// 	imageData, _, err := image.Decode(imageFile)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	buff := new(bytes.Buffer)
	// 	if err = png.Encode(buff, imageData); err != nil {
	// 		panic(err)
	// 	}
	// 	picturesBytes = append(picturesBytes, buff.Bytes())
	// }

	// pictures := []model.Picture{}
	// todo: karna modelnya diganti biar pake `string` untuk nerima base64,
	//       jadinya code dibawah gk bisa dipake
	//
	// for _, value := range picturesBytes {
	// 	pic := model.Picture{
	// 		Data: value,
	// 	}
	// 	pictures = append(pictures, pic)
	// }

	// memory := model.Memory{
	// 	Description:  payload.Description,
	// 	UserID:       payload.UserId,
	// 	MemoriesTags: payload.Tags,
	// 	Pictures:     pictures,
	// }

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
	if err := db.Delete(&model.Memory{}, memoryId).Error; err != nil {
		response := utils.BaseResponse{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.StatusCode = http.StatusNotFound
		}

		return utils.SendResponse(c, response)
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
