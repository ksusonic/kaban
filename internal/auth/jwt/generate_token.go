package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/ksusonic/kanban/internal/models"
)

func (a *Auth) GenerateJWTToken(userID int) (*models.JWTToken, error) {
	now := time.Now()
	expiresAt := now.Add(a.tokenTTL)

	claims := userClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Issuer:    issuer,
		},
		UserID: userID,
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := t.SignedString(a.key)
	if err != nil {
		return nil, err
	}

	return &models.JWTToken{
		Token:   signedToken,
		Expires: expiresAt,
	}, nil
}
