package boards

import (
	"context"
)

const boardDeleteQuery = `
	update boards 
	set 
	    deleted_at = now(),
	    updated_at = now()
	where 
	    id = $1
		and deleted_at is null
`

func (r *Repository) BoardDelete(ctx context.Context, id int) error {
	_, err := r.db.Conn(ctx).Exec(
		ctx,
		boardDeleteQuery,
		id,
	)

	return err
}
