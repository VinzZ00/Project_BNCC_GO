package utils

import (
	"Project_BNCC_GO/config"
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
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

	claims := token.Claims.(*config.JwtClaim)
	return *claims, nil
}

func GetEchoJwtConfig() echojwt.Config {
	return echojwt.Config{
		SigningKey:  config.JWT_KEY,
		TokenLookup: "cookie:Token",
		ErrorHandler: func(c echo.Context, err error) error {
			return SendResponse(c, BaseResponse{
				StatusCode: http.StatusUnauthorized,
				Message:    "You are not authorized",
			})
		},
		ParseTokenFunc: func(c echo.Context, auth string) (interface{}, error) {
			token, err := jwt.ParseWithClaims(auth, &config.JwtClaim{}, func(t *jwt.Token) (interface{}, error) {
				return config.JWT_KEY, nil
			})

			if err != nil {
				return nil, err
			}
			if !token.Valid {
				return nil, errors.New("invalid token")
			}

			return token, nil
		},
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(config.JwtClaim)
		},
	}
}
