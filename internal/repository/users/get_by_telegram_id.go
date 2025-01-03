package users

import (
	"context"

	"github.com/ksusonic/kanban/internal/models"
)

func (r *Repository) GetByTelegramID(ctx context.Context, telegramID int64) (*models.User, error) {
	rows, err := r.db.Conn(ctx).Query(
		ctx,
		`
		select 
		    id,
	    	username,
	    	tg.telegram_id, 
	    	first_name, 
			last_name,
			avatar_url
		from 
			users left join telegram_users tg on users.id = tg.user_id
		where 
			telegram_id = $1`,
		telegramID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var user models.User
	if rows.Next() {
		err = rows.Scan(&user.ID, &user.Username, &user.TelegramID, &user.FirstName, &user.LastName, &user.AvatarURL)
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
