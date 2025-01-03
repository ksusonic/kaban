//go:generate mockgen -source deps.go -destination deps_mock.go -package auth
package auth

import (
	"context"

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
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByTelegramID(ctx context.Context, telegramID int64) (*models.User, error)
	GetUserIDByTelegramID(ctx context.Context, telegramID int64) (int, error)
}
