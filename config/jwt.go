package config

import "github.com/golang-jwt/jwt/v4"

var JWT_KEY []byte = []byte("BrownMinumKopi")

type JwtClaim struct {
	UserID uint `json:"userId"`
	jwt.RegisteredClaims
}
