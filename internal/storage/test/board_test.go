//go:build integration

package test

import (
	"context"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ksusonic/kanban/internal/models"
	"github.com/ksusonic/kanban/internal/storage/postgres"
)

func (suite *IntegrationTestSuite) TestBoard() {
	t := suite.T()

	ctx, err := suite.repo.TransactionContext(context.Background())
	require.NoError(t, err)

	defer func(db *postgres.DB, ctx context.Context) {
		if err = db.Rollback(ctx); err != nil {
			panic(err)
		}
	}(suite.repo.DB, ctx)

	userID, err := suite.repo.UserRepo().AddTelegramUser(ctx, "testman", 123, "tester", nil)
	require.NoError(t, err)

	expectedBoard := models.Board{
		Name:    "Test",
		Slug:    "integration-test",
		OwnerID: userID,
	}

	boardID, err := suite.repo.BoardRepo().BoardAdd(
		ctx,
		expectedBoard.Name,
		expectedBoard.Slug,
		userID,
	)
	require.NoError(t, err)
	expectedBoard.ID = boardID

	actualBoards, err := suite.repo.BoardRepo().BoardsGetAvailable(
		ctx,
		userID,
	)
	require.NoError(t, err)
	assert.Len(t, actualBoards, 1)
	expectedBoard.CreatedAt = actualBoards[0].CreatedAt
	assert.Equal(t, []models.Board{expectedBoard}, actualBoards)

	// double add
	_, err = suite.repo.BoardRepo().BoardAdd(
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
		testGuestID, err = suite.repo.UserRepo().AddTelegramUser(ctx, "test-guest", 124, "tester", nil)
		require.NoError(t, err)

		availableBoards, err = suite.repo.BoardRepo().BoardsGetAvailable(ctx, testGuestID)
		require.NoError(t, err)
		assert.Empty(t, availableBoards)

		// try to double add board
		_, err = suite.repo.BoardRepo().BoardAdd(
			ctx,
			expectedBoard.Name,
			expectedBoard.Slug,
			userID+1,
		)
		require.ErrorIs(t, err, models.ErrAlreadyExists)

		err = suite.repo.BoardMembersRepo().MembersAdd(ctx, boardID, testGuestID, models.AccessLevelRO)
		require.NoError(t, err)

		availableBoards, err = suite.repo.BoardRepo().BoardsGetAvailable(ctx, testGuestID)
		require.NoError(t, err)
		assert.Len(t, availableBoards, 1)
		assert.Equal(t, expectedBoard, availableBoards[0])

		err = suite.repo.BoardMembersRepo().MembersDelete(ctx, boardID, testGuestID)
		require.NoError(t, err)

		availableBoards, err = suite.repo.BoardRepo().BoardsGetAvailable(ctx, testGuestID)
		require.NoError(t, err)
		assert.Empty(t, availableBoards)
	}
}
