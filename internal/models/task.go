package models

import "time"

type Task struct {
	ID          int          `json:"id" db:"id"`
	BoardID     int          `json:"board_id" db:"board_id"`
	AuthorID    int          `json:"author_id" db:"author_id"`
	AssigneeID  int          `json:"assignee_id" db:"assignee_id"`
	Title       string       `json:"title" db:"title"`
	Description string       `json:"description" db:"description"`
	Priority    TaskPriority `json:"task_priority" db:"task_priority"`
	DueDate     *time.Time   `json:"due_date" db:"due_date"`
}

type TaskPriority int

const (
	TaskPriorityLow TaskPriority = iota
	TaskPriorityMedium
	TaskPriorityHigh
)
