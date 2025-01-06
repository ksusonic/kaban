package members

import (
	"context"

	"github.com/ksusonic/kanban/internal/models"
)

const membersAddQuery = `
	insert into board_members (board_id, user_id, access_level)
	values ($1, $2, $3)
	on conflict (board_id, user_id)
		do update set access_level = $3,
					  added_at     = now(),
					  updated_at   = now(),
					  deleted_at   = null`

func (r *Repository) MembersAdd(
	ctx context.Context,
	boardID int,
	userID int,
	accessLevel models.AccessLevel,
) error {
	_, err := r.db.Conn(ctx).Exec(
		ctx,
		membersAddQuery,
		boardID,
		userID,
		accessLevel,
	)

	return err
}
