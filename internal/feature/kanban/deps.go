package kanban

import (
	"context"

	"github.com/ksusonic/kanban/internal/models"
)

type boardRepo interface {
	BoardsGetAvailable(
		ctx context.Context,
		userID int,
	) ([]models.Board, error)
	BoardAdd(
		ctx context.Context,
		name, slug string,
		ownerID int,
	) (int, error)
	BoardDelete(
		ctx context.Context,
		id int,
	) error
}

type boardMemberRepo interface {
	BoardAddMember(
		ctx context.Context,
		boardID int,
		userID int,
		accessLevel models.AccessLevel,
	) error
	BoardDeleteMember(
		ctx context.Context,
		boardID int,
		userID int,
	) error
}
