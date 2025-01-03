//go:build integration

package users_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"github.com/ksusonic/kanban/internal/repository/postgres"
	"github.com/ksusonic/kanban/internal/repository/users"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/ksusonic/kanban/internal/models"
)

func TestRepository_Insert(t *testing.T) {
	ctx := context.Background()

	err := godotenv.Load()
	require.NoError(t, err)

	poolCfg, err := pgxpool.ParseConfig("")
	require.NoError(t, err)

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	require.NoError(t, err)
	defer pool.Close()

	db, closeDB, err := postgres.NewDB(ctx, slog.Default())
	require.NoError(t, err)
	defer closeDB()

	repo := users.NewRepository(db)

	telegramID := int64(123456789)
	user := &models.User{
		TelegramID: &telegramID,
		Username:   "testuser",
		FirstName:  "Test",
	}

	txCtx, err := db.TransactionContext(ctx)
	require.NoError(t, err)

	defer func(db *postgres.DB, ctx context.Context) {
		if err = db.Rollback(ctx); err != nil {
			panic(err)
		}
	}(db, txCtx)

	userID, err := repo.AddTelegramUser(
		txCtx,
		user.Username,
		telegramID,
		user.FirstName,
		user.AvatarURL,
	)
	require.NoError(t, err)

	actual, err := repo.GetByTelegramID(txCtx, telegramID)
	require.NoError(t, err)

	user.ID = actual.ID
	assert.Equal(t, user, actual)

	actualUserID, err := repo.GetUserIDByTelegramID(txCtx, telegramID)
	require.NoError(t, err)
	assert.Equal(t, userID, actualUserID)

	actual, err = repo.GetByID(txCtx, user.ID)
	require.NoError(t, err)

	assert.Equal(t, user, actual)
}
