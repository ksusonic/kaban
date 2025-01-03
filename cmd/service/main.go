package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/ksusonic/kanban/internal/controller"
	"github.com/ksusonic/kanban/internal/controller/auth"
	"github.com/ksusonic/kanban/internal/logger"
	"github.com/ksusonic/kanban/internal/models"
	"github.com/ksusonic/kanban/internal/repository"
)

const (
	ReadTimeout             = 1 * time.Second
	WriteTimeout            = 1 * time.Second
	IdleTimeout             = 30 * time.Second
	ReadHeaderTimeout       = 2 * time.Second
	GracefulShutdownTimeout = 5 * time.Second
)

func main() {
	ctx := context.Background()

	var (
		debugFlag = flag.Bool("debug", false, "enable debug logging")
		httpAddr  = flag.String("http", "localhost:8080", "HTTP service address")
	)

	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	log := logger.New(*debugFlag)

	repo, repoClose, err := repository.NewRepository(ctx, log)
	if err != nil {
		log.Error("init pgx pool", "error", err)
		os.Exit(1)
	}

	defer repoClose()

	router := controller.NewRouter(log, *debugFlag)
	router.LoadHTMLGlob("templates/*.tmpl")

	{
		botCfg := models.BotCfg{
			Name:  os.Getenv("BOT_NAME"),
			Token: os.Getenv("BOT_TOKEN"),
		}
		authCtrl := auth.NewController(repo.UserRepo(), botCfg, log)
		authGroup := router.Group("/auth")

		authGroup.GET("/", authCtrl.Page) // TODO: @sonanted - render page on frontend
		authGroup.GET("/tg-callback", authCtrl.TelegramCallback)
	}

	server := &http.Server{
		Addr:              *httpAddr,
		Handler:           router.Handler(),
		ReadTimeout:       ReadTimeout,
		ReadHeaderTimeout: ReadHeaderTimeout,
		WriteTimeout:      WriteTimeout,
		IdleTimeout:       IdleTimeout,
	}

	run(ctx, server, log)
}

func run(ctx context.Context, srv *http.Server, log logger.Logger) {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log.InfoContext(ctx, "starting server", "addr", srv.Addr)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.ErrorContext(ctx, "serving", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	stop()
	log.InfoContext(ctx, "shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), GracefulShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.ErrorContext(ctx, "force shutdown", err)
	}

	log.InfoContext(ctx, "server stopped")
}
