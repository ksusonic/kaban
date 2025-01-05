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

	_ = godotenv.Load(".env", "../../../.env")

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
	testAvatar := "testtestavatar"
	expectedUser := &models.User{
		Username:   "testuser",
		FirstName:  "Test",
		TelegramID: &telegramID,
		AvatarURL:  &testAvatar,
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
		expectedUser.Username,
		telegramID,
		expectedUser.FirstName,
		expectedUser.AvatarURL,
	)
	require.NoError(t, err)

	expectedUser.ID = userID

	actual, err := repo.GetByID(txCtx, userID)
	require.NoError(t, err)

	assert.Equal(t, expectedUser, actual)

	actual, err = repo.GetByTelegramID(txCtx, telegramID)
	require.NoError(t, err)

	assert.Equal(t, expectedUser, actual)
}
