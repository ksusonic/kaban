package tasks

import (
	"context"
	"fmt"

	"github.com/ksusonic/kanban/internal/models"
)

const tasksAddQuery = `
	insert into tasks (board_id, author_id, assignee_id, title, description, priority, due_date)
	values ($1, $2, $3, $4, $5, $6, $7)
	returning id`

func (r *Repository) TasksAdd(ctx context.Context, task models.Task) (int, error) {
	rows, err := r.db.Conn(ctx).Query(
		ctx,
		tasksAddQuery,
		task.BoardID,
		task.AuthorID,
		task.AssigneeID,
		task.Title,
		task.Description,
		task.Priority,
		task.DueDate,
	)
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	if !rows.Next() {
		return 0, fmt.Errorf("task %d not added: %w", task.AuthorID, rows.Err())
	}

	var id int
	if err = rows.Scan(&id); err != nil {
		return 0, fmt.Errorf("task %d not added: %w", task.AuthorID, err)
	}

	return id, nil
}
