package controller

import (
	"Project_BNCC_GO/model"
	"bytes"
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

	webData := struct {
		Memoryid     uint        `json:"memoryid,omitempty"`
		DateAdded    time.Time   `json:"dateAdded"`
		DateModified time.Time   `json:"dateModified"`
		Desc         string      `json:"MemoryDesc"`
		UserId       uint        `json:"UserId"`
		PictureId    uint        `json:"pictureId,omitempty"`
		Path         []string    `json:"PicturePath"`
		Tag          []model.Tag `json:"tags"`
	}{}

	if err := c.Bind(&webData); err != nil {
		panic("Error di binding data")
	}

	picturesbyte := [][]byte{}

	for _, value := range webData.Path {
		imagefile, err := os.Open(value)
		if err != nil {
			panic("path invalid")
		}
		imageData, _, err := image.Decode(imagefile)
		buff := new(bytes.Buffer)
		err = png.Encode(buff, imageData)

		if err != nil {
			panic(err)
		}

		picturesbyte = append(picturesbyte, buff.Bytes())
	}

	pictures := []model.Picture{}

	for _, value := range picturesbyte {
		pic := model.Picture{
			Picture: value,
		}

		pictures = append(pictures, pic)

	}

	memory := model.Memory{
		Memoryid:     webData.Memoryid,
		DateAdded:    webData.DateAdded,
		DateModified: webData.DateModified,
		Desc:         webData.Desc,
		UserId:       webData.UserId,
		Tag:          webData.Tag,
		Picture:      pictures,
	}
	result := db.Create(&memory)

	fmt.Println(result)

	return c.JSON(http.StatusOK, "Insert Successful")
}

func init() {
	if database, err := model.GetDB(); err == nil {
		db = database
	} else {
		panic(err)
	}

}
