package auth

import (
	"github.com/ksusonic/kanban/internal/logger"
	"github.com/ksusonic/kanban/internal/models"
)

type Controller struct {
	userRepo   userRepo
	authModule authModule
	log        logger.Logger
	botCfg     models.BotCfg
}

func NewController(
	userRepo userRepo,
	authModule authModule,
	log logger.Logger,
	botCfg models.BotCfg,
) *Controller {
	return &Controller{
		userRepo:   userRepo,
		authModule: authModule,
		botCfg:     botCfg,
		log:        log,
	}
}
