//go:generate mockgen -source deps.go -destination deps_mock.go -package board
package board

import (
	"context"

	"github.com/ksusonic/kanban/internal/models"
)

type Feature interface {
	AvailableBoards(ctx context.Context, userID int) ([]models.Board, error)
	GetBoardBySlug(ctx context.Context, userID int, boardSlug string) (*models.Board, error)
}
