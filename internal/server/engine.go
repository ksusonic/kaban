package server

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	auth2 "github.com/ksusonic/kanban/internal/auth"
	boardFeature "github.com/ksusonic/kanban/internal/feature/board"
	"github.com/ksusonic/kanban/internal/models"
	servicemw "github.com/ksusonic/kanban/internal/server/middleware"
	"github.com/ksusonic/kanban/internal/server/routes/auth"
	"github.com/ksusonic/kanban/internal/server/routes/board"
)

func BuildEngine(
	repo Repo,
	log *slog.Logger,
	isDebug bool,
) http.Handler {
	if !isDebug {
		gin.SetMode(gin.ReleaseMode)
	}

	botCfg := models.BotCfg{
		Name:  os.Getenv("BOT_NAME"),
		Token: os.Getenv("BOT_TOKEN"),
	}

	authProvider := auth2.NewAuth(
		os.Getenv("SECRET_KEY"),
		botCfg.Token,
		(time.Hour*24)*7, //nolint:mnd // 7 days TTL
	)

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

		authCtrl := auth.NewController(repo.UserRepo(), authProvider, log, botCfg)

		authGroup.GET("/", authCtrl.Page) // TODO: @sonanted - render page on frontend
		authGroup.GET("/tg-callback", authCtrl.TelegramCallback)
	}
	{
		boardGroup := engine.Group("/boards")
		boardGroup.Use(servicemw.AuthRequired(authProvider))

		boardCtrl := board.NewController(
			boardFeature.New(
				repo.BoardRepo(),
				repo.BoardMembersRepo(),
			),
			log,
		)

		boardGroup.GET("/", boardCtrl.AvailableBoards)
		boardGroup.GET("/:slug", boardCtrl.GetBoardBySlug)

		// to be continued...
	}

	return engine
}
