package controller

import (
	"Project_BNCC_GO/config"
	"Project_BNCC_GO/model"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func Login(c echo.Context) error {
	//Login
	LoginWebData := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	if err := c.Bind(&LoginWebData); err != nil {
		panic(err)
	}
	user := model.User{
		Email: LoginWebData.Email,
	}
	result := db.Where("email = ?", user.Email).Take(&user)
	if err := result.Error; err != nil {
		panic(err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(LoginWebData.Password)); err != nil {
		panic(err)
	}

	//Generate token
	expTime := time.Now().Add(time.Minute * 15)
	claims := config.JwtClaim{
		UserId: user.Userid,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "PROJECT_BNCC_GO",
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}
	Algo := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	Token, err := Algo.SignedString(config.JWT_KEY)
	if err != nil {
		panic(err)
	}

	// Create Cookienya
	cookie := new(http.Cookie)
	cookie.Name = "Token"
	cookie.Value = Token
	cookie.Path = "/"
	cookie.HttpOnly = true

	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, struct {
		Message string
	}{
		Message: "Sucessfully Logged In",
	})
}

func SignUP(c echo.Context) error {
	signupUser := model.User{}

	c.Bind(&signupUser)

	log.Println(signupUser)
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(signupUser.Password), bcrypt.DefaultCost)

	signupUser.Password = string(hashPassword)

	result := db.Create(&signupUser)

	if result.Error != nil {
		panic(result.Error)
	}

	return c.JSON(http.StatusOK, struct {
		Message string
	}{
		Message: "User with ID " + strconv.FormatUint(uint64(signupUser.Userid), 10) + "is Created",
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
