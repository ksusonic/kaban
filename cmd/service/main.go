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

	"github.com/ksusonic/kanban/internal/controller"
	userctrl "github.com/ksusonic/kanban/internal/controller/user"

	"github.com/ksusonic/kanban/internal/logger"
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
	var (
		ctx = context.Background()
		log = logger.New()

		debugFlag = flag.Bool("debug", false, "enable debug logging")
		httpAddr  = flag.String("http", "localhost:8080", "HTTP service address")
	)

	flag.Parse()

	repo, repoClose, err := repository.NewRepository(ctx, log)
	if err != nil {
		log.ErrorContext(ctx, "init pgx pool", "error", err)
		os.Exit(1)
	}

	defer repoClose()

	router := controller.NewRouter(log, *debugFlag)

	{
		apiGroup := router.AddGroup("/api")
		{
			ctrl := userctrl.NewController(repo.UserRepo(), log)
			userGroup := apiGroup.Group("/user")
			userGroup.POST("/login", ctrl.Login)
		}
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
