//go:build integration

package test

import (
	"context"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ksusonic/kanban/internal/models"
	"github.com/ksusonic/kanban/internal/storage/postgres"
)

func (suite *IntegrationTestSuite) TestTask() {
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

	boardID, err := suite.repo.BoardRepo().BoardAdd(
		ctx,
		"test task",
		"TEST-TASK",
		userID,
	)
	require.NoError(t, err)

	actualTasks, err := suite.repo.BoardTasksRepo().TasksGetByBoard(ctx, boardID)
	require.NoError(t, err)
	require.Len(t, actualTasks, 0)

	expectedTask := models.Task{
		BoardID:     boardID,
		AuthorID:    userID,
		AssigneeID:  userID,
		Title:       "Hello, i am test task",
		Description: "no description yet",
		Priority:    models.TaskPriorityHigh,
		DueDate:     nil,
	}

	taskID, err := suite.repo.BoardTasksRepo().TasksAdd(ctx, expectedTask)
	require.NoError(t, err)
	expectedTask.ID = taskID

	actualTasks, err = suite.repo.BoardTasksRepo().TasksGetByBoard(ctx, boardID)
	require.NoError(t, err)
	require.Len(t, actualTasks, 1)

	assert.Equal(t, expectedTask, actualTasks[0])

	actualTask, err := suite.repo.BoardTasksRepo().TasksGetByID(ctx, taskID)
	require.NoError(t, err)
	require.NotNil(t, actualTask)
	assert.Equal(t, expectedTask, *actualTask)

	// DELETE flow
	err = suite.repo.BoardTasksRepo().TasksDelete(ctx, taskID)
	require.NoError(t, err)

	actualTasks, err = suite.repo.BoardTasksRepo().TasksGetByBoard(ctx, boardID)
	require.NoError(t, err)
	require.Len(t, actualTasks, 0)

	_, err = suite.repo.BoardTasksRepo().TasksGetByID(ctx, taskID)
	require.ErrorIs(t, err, models.ErrNotFound)
}
