package tasks

import (
	"context"

	"github.com/ksusonic/kanban/internal/models"
)

const tasksAddQuery = `
	insert into tasks (board_id, author_id, assignee_id, title, description, priority, due_date)
	values ($1, $2, $3, $4, $5, $6, $7)`

func (r *Repository) TasksAdd(ctx context.Context, task models.Task) error {
	_, err := r.db.Conn(ctx).Exec(
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

	return err
}
