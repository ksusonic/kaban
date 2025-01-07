package auth

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"

	"github.com/ksusonic/kanban/internal/models"
)

func (a *Auth) CheckToken(tokenString string) (*models.UserIdentity, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&userClaims{},
		func(_ *jwt.Token) (interface{}, error) {
			return a.key, nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("parse jwt token: %w", err)
	} else if claims, ok := token.Claims.(*userClaims); ok {
		return &models.UserIdentity{UserID: claims.UserID}, nil
	}

	return nil, errors.New("parse jwt token: invalid claims")
}
