package members

import (
	"context"
)

const boardMembersDeleteQuery = `
	update board_members 
	set 
		deleted_at = now()
	where 
	    board_id = $1 
	    and user_id = $2 
	    and deleted_at is null`

func (r *Repository) BoardMembersDelete(
	ctx context.Context,
	boardID int,
	userID int,
) error {
	_, err := r.db.Conn(ctx).Exec(
		ctx,
		boardMembersDeleteQuery,
		boardID,
		userID,
	)

	return err
}
