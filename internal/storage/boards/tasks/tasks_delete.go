package tasks

import (
	"context"
	"fmt"
)

const tasksDeleteQuery = `
	update tasks 
	set deleted_at = now()
	where id = $1 
	  and deleted_at is null
	returning id`

func (r *Repository) TasksDelete(ctx context.Context, id int) error {
	rows, err := r.db.Conn(ctx).Query(
		ctx,
		tasksDeleteQuery,
		id,
	)
	if err != nil {
		return err
	}

	defer rows.Close()

	if !rows.Next() {
		if err = rows.Err(); err != nil {
			return fmt.Errorf("task %d not deleted: %w", id, err)
		}

		return fmt.Errorf("task %d not deleted", id)
	}

	return nil
}
