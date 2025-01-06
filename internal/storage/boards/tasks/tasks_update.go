package tasks

import (
	"context"

	"github.com/ksusonic/kanban/internal/models"
)

const tasksUpdateQuery = `
	update tasks
	set
		board_id = coalesce($2, board_id),
		author_id = coalesce($3, author_id),
		assignee_id = coalesce($4, assignee_id),
		title = coalesce($5, title),
		description = coalesce($6, description),
		priority = coalesce($7, priority),
		due_date = coalesce($8, due_date),
		updated_at = now()
	where id = $1
	returning id, board_id, author_id, assignee_id, title, description, priority, due_date, updated_at`

func (r *Repository) TasksUpdate(ctx context.Context, task models.Task) error {
	_, err := r.db.Conn(ctx).Exec(
		ctx,
		tasksUpdateQuery,
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
