package server

import (
	servicemw "github.com/ksusonic/kanban/internal/server/middleware"
	"github.com/ksusonic/kanban/internal/storage/users"
)

type Repo interface {
	UserRepo() *users.Repository
}

type AuthProvider servicemw.AuthProvider
