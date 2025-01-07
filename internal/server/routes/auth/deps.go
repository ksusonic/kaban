//go:generate mockgen -source deps.go -destination deps_mock.go -package auth
package auth

import (
	"context"
	"net/url"

	"github.com/ksusonic/kanban/internal/models"
)

type userRepo interface {
	AddTelegramUser(
		ctx context.Context,
		username string,
		telegramID int64,
		firstName string,
		avatarURL *string,
	) (int, error)
	GetByTelegramID(ctx context.Context, telegramID int64) (*models.User, error)
}

type authModule interface {
	GenerateJWTToken(userID int) (*models.JWTToken, error)
	ValidateTelegramCallbackData(query url.Values) bool
}
