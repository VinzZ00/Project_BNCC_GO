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

type MemoryResponse struct {
	model.BaseModel
	Description string   `json:"description"`
	UserID      uint     `json:"user_id"`
	Pictures    []string `json:"pictures"`
	Tags        []string `json:"tags"`
}

func AddPicture(c echo.Context) error {
	payload := struct {
		ID   int    `param:"id" validate:"required,number"`
		Data string `json:"base64Data" validate:"required"`
	}{}

	if err := c.Bind(&payload); err != nil {
		return utils.SendResponse(c, utils.BaseResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
	}

	picture := model.Picture{
		Data:     payload.Data,
		MemoryID: uint(payload.ID),
	}
	CheckUser := model.Memory{}
	if err := db.Where("ID = ?", picture.MemoryID).First(&CheckUser).Error; err != nil {
		return utils.SendResponse(c, utils.BaseResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}
	curr_user, err := utils.GetAuthUser(c)
	if err != nil {
		return utils.SendResponse(c, utils.BaseResponse{
			StatusCode: http.StatusNetworkAuthenticationRequired,
			Message:    err.Error(),
		})
	}

	if CheckUser.UserID != curr_user.UserID {
		return utils.SendResponse(c, utils.BaseResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    err.Error(),
		})
	}

	if err := db.Create(&picture).Error; err != nil {
		return utils.SendResponse(c, utils.BaseResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, struct {
		Status  int
		Message string
	}{
		Status:  http.StatusOK,
		Message: fmt.Sprintf("Picture with ID %d has been successfully added", picture.ID),
	})
}

func mapMemoryToResponse(memory model.Memory) MemoryResponse {
	var pictureLinks []string
	var tags []string

	for _, picture := range memory.Pictures {
		link := fmt.Sprintf("/pictures/%d", picture.ID)
		pictureLinks = append(pictureLinks, link)
	}
	for _, memoryTag := range memory.MemoriesTags {
		tagName := memoryTag.Tag.Name
		tags = append(tags, tagName)
	}

	return MemoryResponse{
		BaseModel: model.BaseModel{
			ID:        memory.ID,
			CreatedAt: memory.CreatedAt,
			UpdatedAt: memory.UpdatedAt,
			DeletedAt: memory.DeletedAt,
		},
		Description: memory.Description,
		UserID:      memory.UserID,
		Pictures:    pictureLinks,
		Tags:        tags,
	}
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
	var tags []model.Tag
	for _, val := range payload.Tags {
		var tag model.Tag

		err := db.Where("name = ?", val).First(&tag).Error
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return utils.SendResponse(c, utils.BaseResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    err.Error(),
				})
			}

			// when a tag doesn't exist in the db, we'll add them to the db
			// this allows us to make use of the existing tags, while adding new tag if it never existed.
			tag = model.Tag{
				Name: val,
			}
			if err := db.Create(&tag).Error; err != nil {
				return utils.SendResponse(c, utils.BaseResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    err.Error(),
				})
			}
		}

		tags = append(tags, tag)
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
		return utils.SendResponse(c, utils.BaseResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	return utils.SendResponse(c, utils.BaseResponse{
		StatusCode: http.StatusCreated,
		Message:    "Memory is successfully created",
	})
}

func GetAllMemories(c echo.Context) error {
	memories := []model.Memory{}

	currentUser, _ := utils.GetAuthUser(c)
	if err := db.Where("User_ID = ? ", currentUser.UserID).Preload("Pictures").Preload("MemoriesTags").Preload("MemoriesTags.Tag").Find(&memories).Error; err != nil {
		return utils.SendResponse(c, utils.BaseResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	var mappedResponse []MemoryResponse
	for _, memory := range memories {
		mappedResponse = append(mappedResponse, mapMemoryToResponse(memory))
	}

	return c.JSON(http.StatusOK, &mappedResponse)
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

	if err := db.Where("id = ?", payload.ID).Preload("Pictures").Preload("MemoriesTags").Preload("MemoriesTags.Tag").Find(&memory).Error; err != nil {
		return utils.SendResponse(e, utils.BaseResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	currentUser, _ := utils.GetAuthUser(e)
	if memory.UserID != currentUser.UserID {
		return utils.SendResponse(e, utils.BaseResponse{
			StatusCode: http.StatusForbidden,
			Message:    "You do not own this memory",
		})
	}

	mappedResponse := mapMemoryToResponse(memory)
	return e.JSON(http.StatusOK, &mappedResponse)
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

func GetMemorySortBy(c echo.Context) error {
	payload := struct {
		SortBy string `query:"sort" validate:"oneof='upload_time' 'tags' 'last_edit' ''"`
		Order  string `query:"order" validate:"oneof='asc' 'desc' ''"`
	}{}

	if err := c.Bind(&payload); err != nil {
		return utils.SendResponse(c, utils.BaseResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
	}
	if err := c.Validate(&payload); err != nil {
		return utils.SendResponse(c, utils.BaseResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
	}

	if payload.SortBy != "" && payload.Order == "" {
		payload.Order = "asc"
	}

	currentUser, _ := utils.GetAuthUser(c)
	var memories []model.Memory

	dbQuery := db.Where("user_id = ?", currentUser.UserID).
		Preload("Pictures").
		Preload("MemoriesTags").
		Preload("MemoriesTags.Tag")

	switch payload.SortBy {
	case "upload_time":
		dbQuery = dbQuery.Order("created_at " + payload.Order)
	case "tags":
		// dbQuery = db.Table(
		// 	"(?) as ",
		// 	dbQuery.Table("memory").
		// 		Joins("JOIN memory_tag on memory.id = memory_tag.memory_id").
		// 		Joins("JOIN tag on tag.id = memory_tag.tag_id").
		// 		Order("tag.name "+payload.Order),
		// ).Distinct()
		dbQuery = dbQuery.
			Joins("JOIN memory_tag on memory.id = memory_tag.memory_id").
			Joins("JOIN tag on tag.id = memory_tag.tag_id").
			Order("tag.name " + payload.Order).
			Distinct()
	case "last_edit":
		dbQuery = dbQuery.Order("updated_at " + payload.Order)
	}

	if err := dbQuery.Find(&memories).Error; err != nil {
		return utils.SendResponse(c, utils.BaseResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	var mappedResponse []MemoryResponse
	for _, memory := range memories {
		mappedResponse = append(mappedResponse, mapMemoryToResponse(memory))
	}
	return c.JSON(http.StatusOK, &mappedResponse)
}

func MemoryFilterBy(c echo.Context) error {
	currentUser, _ := utils.GetAuthUser(c)
	Ftag := c.QueryParam("tag")
	Fdesc := c.QueryParam("description")
	memories := []model.Memory{}

	if Ftag == "" && Fdesc == "" {
		GetAllMemories(c)
	} else if Ftag == "" {
		fmt.Println("Fdesc is ", Fdesc)
		fmt.Println("Ftags is ", Ftag)
		if err := db.Joins("JOIN memory_tag on memory.id = memory_tag.memory_id").Joins("JOIN tag on tag.id = memory_tag.tag_id").Where("user_id = ? and description = ?", currentUser.UserID, Fdesc).Preload("Pictures").Preload("MemoriesTags").Preload("MemoriesTags.Tag").Find(&memories).Error; err != nil {
			utils.SendResponse(c, utils.BaseResponse{
				StatusCode: http.StatusPreconditionFailed,
				Message:    err.Error(),
			})
		}
	} else if Fdesc == "" {
		filter := GetTagIdByName(Ftag)
		fmt.Println("Fdesc is ", Fdesc)
		fmt.Println("Ftags is ", Ftag)
		fmt.Println("The tag id is ", filter)
		if err := db.Joins("JOIN memory_tag on memory.id = memory_tag.memory_id").Joins("JOIN tag on tag.id = memory_tag.tag_id").Where("user_id = ? and tag.ID = ?", currentUser.UserID, filter).Preload("Pictures").Preload("MemoriesTags").Preload("MemoriesTags.Tag").Distinct().Find(&memories).Error; err != nil {
			utils.SendResponse(c, utils.BaseResponse{
				StatusCode: http.StatusPreconditionFailed,
				Message:    err.Error(),
			})
		}
	} else {
		filter := GetTagIdByName(Ftag)
		fmt.Println("Fdesc is ", Fdesc)
		fmt.Println("Ftags is ", Ftag)
		fmt.Println("The tag id is ", filter)
		if err := db.Joins("JOIN memory_tag on memory.id = memory_tag.memory_id").Joins("JOIN tag on tag.id = memory_tag.tag_id").Where("user_id = ? and tag.ID = ? and description = ?", currentUser.UserID, filter, Fdesc).Preload("Pictures").Preload("MemoriesTags").Preload("MemoriesTags.Tag").Distinct().Find(&memories).Error; err != nil {
			utils.SendResponse(c, utils.BaseResponse{
				StatusCode: http.StatusPreconditionFailed,
				Message:    err.Error(),
			})
		}
	}

	var mappedResponse []MemoryResponse
	for _, memory := range memories {
		mappedResponse = append(mappedResponse, mapMemoryToResponse(memory))
	}
	return c.JSON(http.StatusOK, &mappedResponse)
}

func GetTagbyTagID(tagId uint) (Tag model.Tag) {
	Tag.ID = tagId
	if err := db.Find(&Tag).Error; err != nil {
		fmt.Println("error")
	}
	return
}

func GetAllTags(c echo.Context) error {
	tags := model.Tag{}

	db.Find(&tags)

	return c.JSON(http.StatusOK, &tags)
}

func GetTagIdByName(tagName string) (TagId uint) {
	tag := model.Tag{}
	if err := db.Where("name = ?", tagName).Find(&tag).Error; err != nil {
		panic(err)
	}
	TagId = tag.ID
	return
}

func init() {
	if database, err := model.GetDB(); err == nil {
		db = database
	} else {
		panic(err)
	}
}
