////go:build integration

package boards_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ksusonic/kanban/internal/models"
	"github.com/ksusonic/kanban/internal/storage/boards"
	"github.com/ksusonic/kanban/internal/storage/boards/members"
	"github.com/ksusonic/kanban/internal/storage/postgres"
	"github.com/ksusonic/kanban/internal/storage/users"
)

func TestRepository_Board(t *testing.T) {
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

	userRepo := users.NewRepository(db)
	repo := boards.NewRepository(db)

	ctx, err := db.TransactionContext(context.Background())
	require.NoError(t, err)

	defer func(db *postgres.DB, ctx context.Context) {
		if err = db.Rollback(ctx); err != nil {
			panic(err)
		}
	}(db, ctx)

	userID, err := userRepo.AddTelegramUser(ctx, "testman", 123, "tester", nil)
	require.NoError(t, err)

	expectedBoard := models.Board{
		Name:    "Test",
		Slug:    "integration-test",
		OwnerID: userID,
	}

	boardID, err := repo.BoardAdd(
		ctx,
		expectedBoard.Name,
		expectedBoard.Slug,
		userID,
	)
	require.NoError(t, err)
	expectedBoard.ID = boardID

	actualBoards, err := repo.BoardsGetAvailable(
		ctx,
		userID,
	)
	require.NoError(t, err)
	assert.Len(t, actualBoards, 1)
	expectedBoard.CreatedAt = actualBoards[0].CreatedAt
	assert.Equal(t, []models.Board{expectedBoard}, actualBoards)

	// double add
	_, err = repo.BoardAdd(
		ctx,
		expectedBoard.Name,
		expectedBoard.Slug,
		userID+1,
	)
	require.ErrorIs(t, err, models.ErrAlreadyExists)

	// board members
	{
		var (
			testGuestID     int
			availableBoards []models.Board
		)
		testGuestID, err = userRepo.AddTelegramUser(ctx, "test-guest", 124, "tester", nil)
		require.NoError(t, err)

		availableBoards, err = repo.BoardsGetAvailable(ctx, testGuestID)
		require.NoError(t, err)
		assert.Empty(t, availableBoards)

		// try to double add board
		_, err = repo.BoardAdd(
			ctx,
			expectedBoard.Name,
			expectedBoard.Slug,
			userID+1,
		)
		require.ErrorIs(t, err, models.ErrAlreadyExists)

		memberRepo := members.NewRepository(db)
		err = memberRepo.BoardAddMember(ctx, boardID, testGuestID, models.AccessLevelRO)
		require.NoError(t, err)

		availableBoards, err = repo.BoardsGetAvailable(ctx, testGuestID)
		require.NoError(t, err)
		assert.Len(t, availableBoards, 1)
		assert.Equal(t, expectedBoard, availableBoards[0])

		err = memberRepo.BoardMembersDelete(ctx, boardID, testGuestID)
		require.NoError(t, err)

		availableBoards, err = repo.BoardsGetAvailable(ctx, testGuestID)
		require.NoError(t, err)
		assert.Empty(t, availableBoards)
	}
}
