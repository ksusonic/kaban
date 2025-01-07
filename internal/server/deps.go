package server

import (
	servicemw "github.com/ksusonic/kanban/internal/server/middleware"
	"github.com/ksusonic/kanban/internal/storage/boards"
	"github.com/ksusonic/kanban/internal/storage/boards/members"
	"github.com/ksusonic/kanban/internal/storage/boards/tasks"
	"github.com/ksusonic/kanban/internal/storage/users"
)

type Repo interface {
	UserRepo() *users.Repository
	BoardRepo() *boards.Repository
	BoardMembersRepo() *members.Repository
	BoardTasksRepo() *tasks.Repository
}

type AuthProvider servicemw.AuthProvider
