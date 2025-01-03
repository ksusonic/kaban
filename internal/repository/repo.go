package repository

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ksusonic/kanban/internal/repository/postgres"
	"github.com/ksusonic/kanban/internal/repository/users"
)

type Repository struct {
	userRepo *users.Repository
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
		userRepo: users.NewRepository(db),
	}, closer, err
}

func (r *Repository) UserRepo() *users.Repository {
	return r.userRepo
}
