package tasks

import (
	"context"

	"github.com/ksusonic/kanban/internal/models"
)

const tasksGetByIDQuery = `
	select 
	    id,
	    board_id, 
	    author_id, 
	    assignee_id, 
	    title, 
	    description, 
	    priority, 
	    due_date
	from tasks
	where 
	    id = $1
	    and deleted_at is null`

func (r *Repository) TasksGetByID(ctx context.Context, id int) (*models.Task, error) {
	rows, err := r.db.Conn(ctx).Query(
		ctx,
		tasksGetByIDQuery,
		id,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var task models.Task
	if rows.Next() {
		err = rows.Scan(
			&task.ID,
			&task.BoardID,
			&task.AuthorID,
			&task.AssigneeID,
			&task.Title,
			&task.Description,
			&task.Priority,
			&task.DueDate,
		)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, models.ErrNotFound
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &task, err
}
