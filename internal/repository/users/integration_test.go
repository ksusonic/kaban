//go:build integration

package users_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/joho/godotenv"
	"github.com/ksusonic/kanban/internal/repository/postgres"
	"github.com/ksusonic/kanban/internal/repository/users"
	"github.com/stretchr/testify/assert"

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
		TelegramID: telegramID,
		Username:   "testuser",
		FirstName:  "Test",
		LastName:   "User",
	}

	txCtx, err := db.TransactionContext(ctx)
	require.NoError(t, err)

	defer func(db *postgres.DB, ctx context.Context) {
		if err = db.Rollback(ctx); err != nil {
			panic(err)
		}
	}(db, txCtx)

	id, err := repo.Add(txCtx, user)
	require.NoError(t, err)

	actual, err := repo.GetByID(txCtx, id)
	require.NoError(t, err)

	assert.Equal(t, user, actual)

	actual, err = repo.GetByTelegramID(txCtx, telegramID)
	require.NoError(t, err)

	assert.Equal(t, user, actual)
}
