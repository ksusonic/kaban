package board

import (
	"github.com/ksusonic/kanban/internal/logger"
)

type Controller struct {
	feature Feature
	log     logger.Logger
}

func NewController(
	feature Feature,
	log logger.Logger,
) *Controller {
	return &Controller{
		feature: feature,
		log:     log,
	}
}
