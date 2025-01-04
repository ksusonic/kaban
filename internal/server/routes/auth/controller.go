package auth

import (
	"os"

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
) *Controller {
	botCfg := models.BotCfg{
		Name:  os.Getenv("BOT_NAME"),
		Token: os.Getenv("BOT_TOKEN"),
	}

	return &Controller{
		userRepo:   userRepo,
		authModule: authModule,
		botCfg:     botCfg,
		log:        log,
	}
}
