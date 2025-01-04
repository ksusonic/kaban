package server

import (
	"github.com/ksusonic/kanban/internal/repository/users"
	servicemw "github.com/ksusonic/kanban/internal/server/middleware"
)

type Repo interface {
	UserRepo() *users.Repository
}

type AuthProvider servicemw.AuthProvider
