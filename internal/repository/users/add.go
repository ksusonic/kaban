package users

import (
	"context"
	"fmt"

	"github.com/ksusonic/kanban/internal/models"
)

func (r *Repository) Add(ctx context.Context, user *models.User) (int, error) {
	rows, err := r.db.Conn(ctx).Query(
		ctx,
		`
		insert into users (telegram_id, username, first_name, last_name)
			values ($1, $2, $3, $4)
		returning id`,
		user.TelegramID,
		user.Username,
		user.FirstName,
		user.LastName,
	)
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	if !rows.Next() {
		return 0, fmt.Errorf("user not added: %v", user)
	}

	var id int
	return id, rows.Scan(&id)
}
