package user

import (
	"github.com/ksusonic/kanban/internal/logger"
)

type Controller struct {
	userRepo userRepo
	log      logger.Logger
}

func NewController(userRepo userRepo, log logger.Logger) *Controller {
	return &Controller{
		userRepo: userRepo,
		log:      log,
	}
}
