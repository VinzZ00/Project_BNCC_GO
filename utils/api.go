package utils

import "github.com/labstack/echo/v4"

type BaseResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

// Sends standardized JSON response for the app
func SendResponse(context echo.Context, response BaseResponse) error {
	return context.JSON(response.StatusCode, response)
}
