//go:build integration

package users_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ksusonic/kanban/internal/models"
	"github.com/ksusonic/kanban/internal/storage/postgres"
	"github.com/ksusonic/kanban/internal/storage/users"
)

func TestRepository(t *testing.T) {
	err := godotenv.Load("../../../.env")
	require.NoError(t, err)

	poolCfg, err := pgxpool.ParseConfig("")
	require.NoError(t, err)

	pool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	require.NoError(t, err)
	defer pool.Close()

	db, closeDB, err := postgres.NewDB(context.Background(), slog.Default())
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

	ctx, err := db.TransactionContext(context.Background())
	require.NoError(t, err)

	defer func(db *postgres.DB, ctx context.Context) {
		if err = db.Rollback(ctx); err != nil {
			panic(err)
		}
	}(db, ctx)

	userID, err := repo.AddTelegramUser(
		ctx,
		expectedUser.Username,
		telegramID,
		expectedUser.FirstName,
		expectedUser.AvatarURL,
	)
	require.NoError(t, err)

	expectedUser.ID = userID

	actual, err := repo.GetByID(ctx, userID)
	require.NoError(t, err)

	assert.Equal(t, expectedUser, actual)

	actual, err = repo.GetByTelegramID(ctx, telegramID)
	require.NoError(t, err)

	assert.Equal(t, expectedUser, actual)
}
