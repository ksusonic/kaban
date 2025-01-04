package middleware

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

func Sloggin(logger *slog.Logger) gin.HandlerFunc {
	return sloggin.NewWithConfig(logger, sloggin.Config{
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
	})
}
