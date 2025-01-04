package server

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/ksusonic/kanban/internal/auth/jwt"
	servicemw "github.com/ksusonic/kanban/internal/server/middleware"
	"github.com/ksusonic/kanban/internal/server/routes/auth"
)

func BuildEngine(
	repo Repo,
	log *slog.Logger,
	isDebug bool,
) http.Handler {
	if !isDebug {
		gin.SetMode(gin.ReleaseMode)
	}

	authProvider := jwt.NewAuth(os.Getenv("SECRET_KEY"))

	engine := gin.New()
	engine.LoadHTMLGlob("templates/*.tmpl")

	engine.Use(servicemw.Sloggin(log))
	engine.Use(gin.Recovery())

	engine.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello world! Скоро здесь будет 🐗")
	})
	engine.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	{
		authGroup := engine.Group("/auth")
		authCtrl := auth.NewController(repo.UserRepo(), authProvider, log)

		authGroup.GET("/", authCtrl.Page) // TODO: @sonanted - render page on frontend
		authGroup.GET("/tg-callback", authCtrl.TelegramCallback)
	}
	{
		boardGroup := engine.Group("/board")
		boardGroup.Use(servicemw.AuthRequired(authProvider))

		// to be continued...
	}

	return engine
}
