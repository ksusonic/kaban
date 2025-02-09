package storage

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ksusonic/kanban/internal/storage/boards"
	"github.com/ksusonic/kanban/internal/storage/boards/members"
	"github.com/ksusonic/kanban/internal/storage/boards/tasks"
	"github.com/ksusonic/kanban/internal/storage/postgres"
	"github.com/ksusonic/kanban/internal/storage/users"
)

type Repository struct {
	*postgres.DB
	userRepo     *users.Repository
	boardRepo    *boards.Repository
	boardMembers *members.Repository
	boardTasks   *tasks.Repository
}

func NewRepository(
	ctx context.Context,
	log *slog.Logger,
) (*Repository, func(), error) {
	db, closer, err := postgres.NewDB(ctx, log)
	if err != nil {
		return nil, nil, fmt.Errorf("create pgxPool: %w", err)
	}

	return &Repository{
		DB:           db,
		userRepo:     users.NewRepository(db),
		boardRepo:    boards.NewRepository(db),
		boardMembers: members.NewRepository(db),
		boardTasks:   tasks.NewRepository(db),
	}, closer, err
}

func (r *Repository) UserRepo() *users.Repository {
	return r.userRepo
}

func (r *Repository) BoardRepo() *boards.Repository {
	return r.boardRepo
}

func (r *Repository) BoardMembersRepo() *members.Repository {
	return r.boardMembers
}

func (r *Repository) BoardTasksRepo() *tasks.Repository {
	return r.boardTasks
}
