package controller

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

type Router struct {
	engine *gin.Engine
}

func NewRouter(logger *slog.Logger, isDebug bool) *Router {
	if !isDebug {
		gin.SetMode(gin.ReleaseMode)
	}

	slogginConfig := sloggin.Config{
		DefaultLevel:     slog.LevelInfo,
		ClientErrorLevel: slog.LevelWarn,
		ServerErrorLevel: slog.LevelError,

		WithUserAgent:      false,
		WithRequestID:      true,
		WithRequestBody:    false,
		WithRequestHeader:  false,
		WithResponseBody:   false,
		WithResponseHeader: false,
		WithSpanID:         false,
		WithTraceID:        false,
	}

	engine := gin.New()
	engine.Use(sloggin.NewWithConfig(logger, slogginConfig))
	engine.Use(gin.Recovery())

	engine.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return &Router{engine: engine}
}

func (r *Router) Handler() http.Handler {
	return r.engine
}

func (r *Router) AddGroup(relativePath string, handlers ...gin.HandlerFunc) gin.IRouter {
	return r.engine.Group(relativePath, handlers...)
}
