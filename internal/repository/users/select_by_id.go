package users

import (
	"context"

	"github.com/ksusonic/kanban/internal/models"
)

func (r *Repository) SelectByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User

	rows, err := r.db.Conn(ctx).Query(
		ctx,
		`
		select 
			username,
			telegram_id, 
			first_name, 
			last_name,
			avatar_url
		from 
			users 
		where 
			id = $1`,
		id,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&user.Username, &user.TelegramID, &user.FirstName, &user.LastName, &user.AvatarURL)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, models.ErrNotFound
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &user, err
}
