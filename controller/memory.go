package controller

import (
	"Project_BNCC_GO/model"
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
		Description string      `json:"description"`
		UserId      uint        `json:"userId"`
		Paths       []string    `json:"picturePaths"`
		Tags        []model.Tag `json:"tags"`
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
	for _, value := range picturesBytes {
		pic := model.Picture{
			Picture: value,
		}
		pictures = append(pictures, pic)
	}

	currentTime := time.Now()
	memory := model.Memory{
		Desc:         body.Description,
		UserId:       body.UserId,
		Tag:          body.Tags,
		Picture:      pictures,
		DateAdded:    currentTime,
		DateModified: currentTime,
	}
	if err := db.Create(&memory).Error; err != nil {
		panic(err)
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"status":  fmt.Sprint(http.StatusCreated),
		"message": "Memory is successfully created",
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

	memory.Desc = body.Description
	memory.DateModified = time.Now()
	memory.Tag = body.Tags

	if err := db.Save(&memory).Error; err != nil {
		panic(err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status":  fmt.Sprint(http.StatusOK),
		"message": "Memory has been updated",
	})
}

func init() {
	if database, err := model.GetDB(); err == nil {
		db = database
	} else {
		panic(err)
	}
}
