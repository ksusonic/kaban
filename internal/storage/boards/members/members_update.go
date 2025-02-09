package members

import (
	"context"

	"github.com/ksusonic/kanban/internal/models"
)

const membersUpdateQuery = `
	update board_members 
	set
	    access_level = coalesce($1, access_level),
	    updated_at = now()
	where 
	    board_id = $1
		and user_id = $2
		and deleted_at is null
	`

func (r *Repository) MembersUpdate(
	ctx context.Context,
	boardID, userID int,
	accessLevel *models.AccessLevel,
) error {
	_, err := r.db.Conn(ctx).Exec(
		ctx,
		membersUpdateQuery,
		boardID,
		userID,
		accessLevel,
	)

	return err
}
