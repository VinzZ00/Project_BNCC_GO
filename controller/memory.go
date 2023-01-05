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
		Description   string   `json:"description" validate:"require"`
		Base64Picture []string `json:"pictures"`
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
	for _, val := range payload.Base64Picture {
		picture := model.Picture{
			Data: val,
		}
		pictures = append(pictures, picture)
	}

	//Tags Creation
	tags := []model.Tag{}
	for _, val := range payload.Tags {
		checkTag := model.Tag{}

		err := db.Where("name= ?", val).First(&checkTag).Error
		if err != nil {
			tag := model.Tag{
				Name: val,
			}

			if err := db.Create(&tag).Error; err != nil {
				panic(err)
			}
		}
	}

	for _, val := range payload.Tags {
		t := model.Tag{}
		db.Where("name = ? ", val).First(&t)
		tags = append(tags, t)
	}

	//memoryTag
	memoryTags := []model.MemoryTag{}

	for _, val := range tags {
		memoryTag := model.MemoryTag{
			TagID: val.ID,
		}
		memoryTags = append(memoryTags, memoryTag)
	}

	currentUser, _ := utils.GetAuthUser(c)
	memory := model.Memory{
		Description:  payload.Description,
		UserID:       currentUser.UserID,
		Pictures:     pictures,
		MemoriesTags: memoryTags,
	}

	if err := db.Create(&memory).Error; err != nil {
		panic(err)
	}

	return utils.SendResponse(c, utils.BaseResponse{
		StatusCode: http.StatusCreated,
		Message:    "Memory is successfully created",
	})
}

func GetAllMemories(c echo.Context) error {
	memories := []model.Memory{}

	currentUser, _ := utils.GetAuthUser(c)
	fmt.Println("Issued by", currentUser.UserID)

	if err := db.Where("User_ID = ? ", currentUser.UserID).Preload("Pictures").Preload("MemoriesTags").Find(&memories).Error; err != nil {
		return utils.SendResponse(c, utils.BaseResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &memories)
}

func GetAMemories(e echo.Context) error {
	memory := model.Memory{}

	payload := MemoryIDParam{}
	if err := e.Bind(&payload); err != nil {
		return utils.SendResponse(e, utils.BaseResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
	}

	if err := db.Where("id = ?", payload.ID).Preload("Pictures").Preload("MemoriesTags").Find(&memory).Error; err != nil {
		return utils.SendResponse(e, utils.BaseResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	currentUser, _ := utils.GetAuthUser(e)
	fmt.Println("Issued by", currentUser.UserID)

	if memory.UserID != currentUser.UserID {
		return utils.SendResponse(e, utils.BaseResponse{
			StatusCode: http.StatusForbidden,
			Message:    "You do not own this memory",
		})
	}

	return e.JSON(http.StatusOK, &memory)
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
