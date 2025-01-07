package board

import (
	"context"

	"github.com/ksusonic/kanban/internal/models"
)

func (f *Feature) AvailableBoards(ctx context.Context, userID int) ([]models.Board, error) {
	return f.boardRepo.BoardsGetAvailable(ctx, userID)
}
