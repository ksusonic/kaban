package boards

import (
	"context"

	"github.com/ksusonic/kanban/internal/models"
)

const boardsGetBySlugQuery = `
	select 
	    id,
	    name,
	    slug,
	    owner_id,
	    created_at
	from 
	    boards 
	where 
	    slug = $1
		and deleted_at is null`

func (r *Repository) BoardsGetBySlug(ctx context.Context, slug string) (*models.Board, error) {
	rows, err := r.db.Conn(ctx).Query(
		ctx,
		boardsGetBySlugQuery,
		slug,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var board models.Board
	if rows.Next() {
		if err = rows.Scan(
			&board.ID,
			&board.Name,
			&board.Slug,
			&board.OwnerID,
			&board.CreatedAt,
		); err != nil {
			return nil, err
		}
	} else {
		return nil, models.ErrNotFound
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &board, nil
}
