package tasks

import "github.com/ksusonic/kanban/internal/storage/postgres"

type Repository struct {
	db *postgres.DB
}

func NewRepository(db *postgres.DB) *Repository {
	return &Repository{db: db}
}
