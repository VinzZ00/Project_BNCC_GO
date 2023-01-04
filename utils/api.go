package utils

import (
	"Project_BNCC_GO/config"
	"errors"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type BaseResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

// Sends standardized JSON response for the app
func SendResponse(context echo.Context, response BaseResponse) error {
	return context.JSON(response.StatusCode, response)
}

// Grabs the currently authenticated user
// An error will occur if user is unauthorized.
func GetAuthUser(context echo.Context) (config.JwtClaim, error) {
	token := context.Get("user").(*jwt.Token)
	if token == nil {
		return config.JwtClaim{}, errors.New("user is unauthorized")
	}

	return token.Claims.(config.JwtClaim), nil
}
