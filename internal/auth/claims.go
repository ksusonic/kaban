package auth

import "github.com/golang-jwt/jwt/v5"

type userClaims struct {
	jwt.RegisteredClaims
	UserID int `json:"user_id"`
}
