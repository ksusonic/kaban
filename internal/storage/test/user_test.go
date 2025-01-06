//go:build integration

package test

import (
	"context"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ksusonic/kanban/internal/models"
	"github.com/ksusonic/kanban/internal/storage/postgres"
)

func (suite *IntegrationTestSuite) TestUser() {
	t := suite.T()

	telegramID := int64(123456789)
	testAvatar := "testtestavatar"
	expectedUser := &models.User{
		Username:   "testuser",
		FirstName:  "Test",
		TelegramID: &telegramID,
		AvatarURL:  &testAvatar,
	}

	ctx, err := suite.repo.TransactionContext(context.Background())
	require.NoError(t, err)

	defer func(db *postgres.DB, ctx context.Context) {
		if err = db.Rollback(ctx); err != nil {
			panic(err)
		}
	}(suite.repo.DB, ctx)

	userID, err := suite.repo.UserRepo().AddTelegramUser(
		ctx,
		expectedUser.Username,
		telegramID,
		expectedUser.FirstName,
		expectedUser.AvatarURL,
	)
	require.NoError(t, err)

	expectedUser.ID = userID

	actual, err := suite.repo.UserRepo().GetByID(ctx, userID)
	require.NoError(t, err)

	assert.Equal(t, expectedUser, actual)

	actual, err = suite.repo.UserRepo().GetByTelegramID(ctx, telegramID)
	require.NoError(t, err)

	assert.Equal(t, expectedUser, actual)
}
