package users

import (
	"context"

	"github.com/ksusonic/kanban/internal/models"
)

func (r *Repository) GetUserIDByTelegramID(ctx context.Context, telegramID int64) (int, error) {
	rows, err := r.db.Conn(ctx).Query(
		ctx,
		`
		select 
		    id
		from 
			users left join telegram_users tg on users.id = tg.user_id
		where 
			telegram_id = $1`,
		telegramID,
	)
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	var userID int
	if rows.Next() {
		err = rows.Scan(&userID)
		if err != nil {
			return 0, err
		}
	} else {
		return 0, models.ErrNotFound
	}

	if err = rows.Err(); err != nil {
		return 0, err
	}

	return userID, err
}
