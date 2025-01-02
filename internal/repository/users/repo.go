package users

import "github.com/ksusonic/kanban/internal/repository/postgres"

type Repository struct {
	db *postgres.DB
}

func NewRepository(db *postgres.DB) *Repository {
	return &Repository{db: db}
}
