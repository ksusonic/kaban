package tasks

import (
	"context"

	"github.com/ksusonic/kanban/internal/models"
)

const tasksGetByBoardQuery = `
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
	    board_id = $1
	    and deleted_at is null`

func (r *Repository) TasksGetByBoard(ctx context.Context, boardID int) ([]models.Task, error) {
	rows, err := r.db.Conn(ctx).Query(
		ctx,
		tasksGetByBoardQuery,
		boardID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tasks := make([]models.Task, 0, len(rows.RawValues()))

	for rows.Next() {
		var task models.Task

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

		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, err
}
