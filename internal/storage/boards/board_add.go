package boards

import (
	"context"
	"fmt"

	"github.com/ksusonic/kanban/internal/models"
)

const boardAddQuery = `
	insert into boards (name, slug, owner_id) 
	values ($1, $2, $3)
	on conflict (slug) do nothing
	returning id
`

func (r *Repository) BoardAdd(ctx context.Context, name, slug string, ownerID int) (int, error) {
	rows, err := r.db.Conn(ctx).Query(
		ctx,
		boardAddQuery,
		name,
		slug,
		ownerID,
	)
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	if !rows.Next() {
		if err = rows.Err(); err != nil {
			return 0, fmt.Errorf("board %s not added: %w", name, err)
		}

		return 0, models.ErrAlreadyExists
	}

	var id int
	if err = rows.Scan(&id); err != nil {
		return 0, fmt.Errorf("board %s not added: %w", name, err)
	}

	return id, nil
}
