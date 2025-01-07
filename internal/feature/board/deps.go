package board

import (
	"context"

	"github.com/ksusonic/kanban/internal/models"
)

type boardRepo interface {
	BoardsGetAvailable(ctx context.Context, userID int) ([]models.Board, error)
	BoardsGetBySlug(ctx context.Context, slug string) (*models.Board, error)
	BoardAdd(
		ctx context.Context,
		name, slug string,
		ownerID int,
	) (int, error)
	BoardDelete(ctx context.Context, id int) error
}

type boardMemberRepo interface {
	MembersAdd(
		ctx context.Context,
		boardID int,
		userID int,
		accessLevel models.AccessLevel,
	) error
	MembersGet(
		ctx context.Context,
		boardID int,
		userID int,
	) (*models.AccessLevel, error)
	MembersDelete(
		ctx context.Context,
		boardID int,
		userID int,
	) error
}
