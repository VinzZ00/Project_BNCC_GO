package controller

import (
	"Project_BNCC_GO/model"
	"bytes"
	"encoding/json"
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
	result := db.Create(&memory)
	fmt.Println(result)

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)
	res := struct {
		Status  int
		Message string
	}{
		Status:  202,
		Message: "Memory has successfully created",
	}
	return json.NewEncoder(c.Response()).Encode(res)
}

func init() {
	if database, err := model.GetDB(); err == nil {
		db = database
	} else {
		panic(err)
	}
}
