package members

import (
	"context"

	"github.com/ksusonic/kanban/internal/models"
)

const membersGetQuery = `
	select 
	    access_level
	from 
	    board_members
	where 
	  	board_id = $1
	    and user_id = $2
	  	and deleted_at is null`

func (r *Repository) MembersGet(
	ctx context.Context,
	boardID int,
	userID int,
) (*models.AccessLevel, error) {
	rows, err := r.db.Conn(ctx).Query(
		ctx,
		membersGetQuery,
		boardID,
		userID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var accessLevel models.AccessLevel
	if rows.Next() {
		if err = rows.Scan(&accessLevel); err != nil {
			return nil, err
		}
	} else {
		return nil, models.ErrNotFound
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &accessLevel, nil
}
