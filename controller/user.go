package controller

import (
	"Project_BNCC_GO/config"
	"Project_BNCC_GO/model"
	"Project_BNCC_GO/utils"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Login(c echo.Context) error {
	//Login
	payload := struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}{}

	if err := c.Bind(&payload); err != nil {
		return utils.SendResponse(c, utils.BaseResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
	}

	user := model.User{
		Email: payload.Email,
	}
	if err := db.Where("email = ?", user.Email).Take(&user).Error; err != nil {
		response := utils.BaseResponse{
			Message: err.Error(),
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.StatusCode = http.StatusNotFound
		} else {
			response.StatusCode = http.StatusInternalServerError
		}

		return utils.SendResponse(c, response)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		return utils.SendResponse(c, utils.BaseResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Password doesn't match, please try again",
		})
	}

	//Generate token
	expTime := time.Now().Add(time.Minute * 15)
	claims := config.JwtClaim{
		UserId: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "PROJECT_BNCC_GO",
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	algo := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := algo.SignedString(config.JWT_KEY)
	if err != nil {
		return utils.SendResponse(c, utils.BaseResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
	}

	// Create Cookienya
	cookie := new(http.Cookie)
	cookie.Name = "Token"
	cookie.Value = token
	cookie.Path = "/"
	cookie.HttpOnly = true

	c.SetCookie(cookie)

	return utils.SendResponse(c, utils.BaseResponse{
		StatusCode: http.StatusOK,
		Message:    "Successfully logged in",
	})
}

func SignUP(c echo.Context) error {
	payload := struct {
		Username string `json:"username" validate:"required"`
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}{}

	if err := c.Bind(&payload); err != nil {
		return utils.SendResponse(c, utils.BaseResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
	}

	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	newUser := model.User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: string(hashPassword),
	}

	if err := db.Create(&newUser).Error; err != nil {
		return utils.SendResponse(c, utils.BaseResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	return utils.SendResponse(c, utils.BaseResponse{
		StatusCode: http.StatusCreated,
		Message:    fmt.Sprintf("Successfully registered a new user with ID %d", newUser.ID),
	})
}

func init() {
	if db == nil {
		if database, err := model.GetDB(); err == nil {
			db = database
		} else {
			panic(err)
		}
	}
}
