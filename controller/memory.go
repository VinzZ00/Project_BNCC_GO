package controller

import (
	"Project_BNCC_GO/model"
	"Project_BNCC_GO/utils"
	"bytes"
	"fmt"
	"image"
	"image/png"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var db *gorm.DB

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

	currentTime := time.Now()
	memory := model.Memory{
		BaseModel: model.BaseModel{
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
		},
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
		Message: "Memory is successfully created",
	})
}

func UpdateMemory(c echo.Context) error {
	rawId := c.Param("id")
	memoryId, err := strconv.Atoi(rawId)
	if err != nil {
		panic(err)
	}

	body := struct {
		Description string      `json:"description"`
		Tags        []model.Tag `json:"tags"`
	}{}

	if err := c.Bind(&body); err != nil {
		panic(err)
	}

	var memory model.Memory
	if err := db.First(&memory, "memoryid = ?", memoryId).Error; err != nil {
		panic(err)
	}

	// membuat tag yang terupdate baru dari response web
	tags := []model.MemoryTag{}
	for _, tag := range body.Tags {
		tags = append(tags, model.MemoryTag{
			MemoryID: memory.ID,
			TagID:    tag.ID,
		})
	}

	memory.Description = body.Description
	memory.UpdatedAt = time.Now()
	memory.MemoriesTags = tags

	if err := db.Save(&memory).Error; err != nil {
		panic(err)
	}

	return utils.SendResponse(c, utils.BaseResponse{
		StatusCode: http.StatusOK,
		Message: "Memory has been updated",
	})
}

func DeleteMemory(c echo.Context) error {
	rawId := c.Param("id")
	memoryId, err := strconv.Atoi(rawId)
	if err != nil {
		panic(err)
	}

	if err = db.Delete(&model.Memory{}, memoryId).Error; err != nil {
		panic(err)
	}

	return utils.SendResponse(c, utils.BaseResponse{
		StatusCode: http.StatusAccepted,
		Message: fmt.Sprintf("Memory with ID %d has been deleted", memoryId),
	})
}

func init() {
	if database, err := model.GetDB(); err == nil {
		db = database
	} else {
		panic(err)
	}
}
