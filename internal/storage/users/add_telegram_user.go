package users

import (
	"context"
	"fmt"
)

const insertQuery = `
		with new_user as (
			insert into users (username, first_name, avatar_url)
				values ($2, $3, $4)
				returning id)
		insert
		into telegram_users (telegram_id, user_id)
		select $1, id
		from new_user
		returning user_id`

func (r *Repository) AddTelegramUser(
	ctx context.Context,
	username string,
	telegramID int64,
	firstName string,
	avatarURL *string,
) (int, error) {
	rows, err := r.db.Conn(ctx).Query(
		ctx,
		insertQuery,
		telegramID,
		username,
		firstName,
		avatarURL,
	)
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	if !rows.Next() {
		if err = rows.Err(); err != nil {
			return 0, fmt.Errorf("user %s not added: %w", username, err)
		}

		return 0, fmt.Errorf("user %s not added", username)
	}

	var id int
	return id, rows.Scan(&id)
}
