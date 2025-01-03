package auth

import (
	"github.com/ksusonic/kanban/internal/logger"
	"github.com/ksusonic/kanban/internal/models"
)

type Controller struct {
	userRepo userRepo
	log      logger.Logger
	botCfg   models.BotCfg
}

func NewController(userRepo userRepo, botCfg models.BotCfg, log logger.Logger) *Controller {
	return &Controller{
		userRepo: userRepo,
		botCfg:   botCfg,
		log:      log,
	}
}
