package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ksusonic/kanban/internal/repository"
)

func main() {
	var (
		ctx = context.Background()
		log = slog.Default()

		httpAddr = flag.String("http", "localhost:8080", "HTTP service address")
	)

	flag.Parse()

	_, repoClose, err := repository.NewRepository(ctx, log)
	if err != nil {
		log.Error("init pgx pool", "error", err)
		os.Exit(1)
	}

	defer repoClose()

	app := gin.Default()
	app.Use(gin.Recovery())

	if err = app.Run(*httpAddr); err != nil {
		return
	}
}
