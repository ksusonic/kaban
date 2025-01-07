package board

import (
	"context"
	"errors"

	"github.com/ksusonic/kanban/internal/models"
)

func (f *Feature) GetBoardBySlug(
	ctx context.Context,
	userID int,
	boardSlug string,
) (*models.Board, error) {
	board, err := f.boardRepo.BoardsGetBySlug(ctx, boardSlug)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, models.ErrNotFound
		}

		return nil, err
	}

	if board.OwnerID == userID {
		return board, nil
	}

	_, err = f.boardMemberRepo.MembersGet(ctx, board.ID, userID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, models.ErrNoAccess
		}

		return nil, err
	}

	return board, nil
}
