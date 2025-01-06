package boards

import (
	"context"

	"github.com/ksusonic/kanban/internal/models"
)

const boardsGetAvailableQuery = `
	select 
	    id,
	    name,
	    slug,
	    owner_id,
	    created_at
	from 
	    boards 
	where 
	    owner_id = $1
		and deleted_at is null
	union
	select
	    id,
	    name,
	    slug,
	    owner_id,
	    created_at
	from
	    board_members bm left join boards b on b.id = bm.board_id
	where bm.user_id = $1 
	  	and bm.deleted_at is null
`

// BoardsGetAvailable gets all available boards for user: owner or member.
func (r *Repository) BoardsGetAvailable(ctx context.Context, userID int) ([]models.Board, error) {
	rows, err := r.db.Conn(ctx).Query(
		ctx,
		boardsGetAvailableQuery,
		userID,
	)
	if err != nil {
		return nil, err
	}

	boards := make([]models.Board, 0, len(rows.RawValues()))

	for rows.Next() {
		var board models.Board
		if err = rows.Scan(
			&board.ID,
			&board.Name,
			&board.Slug,
			&board.OwnerID,
			&board.CreatedAt,
		); err != nil {
			return nil, err
		}

		boards = append(boards, board)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return boards, nil
}
