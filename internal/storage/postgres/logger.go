package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/tracelog"
)

type pgxLogWrapper struct {
	Logger *slog.Logger
}

func (l *pgxLogWrapper) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
	attrs := make([]slog.Attr, 0, len(data)+1)
	attrs = append(attrs, slog.String("pgx_level", level.String()))
	for k, v := range data {
		attrs = append(attrs, slog.Any(k, v))
	}

	l.Logger.LogAttrs(ctx, slogLevel(level), msg, attrs...)
}

func logLevelFromEnv() tracelog.LogLevel {
	if level := os.Getenv("PGX_LOG_LEVEL"); level != "" {
		l, err := tracelog.LogLevelFromString(level)
		if err != nil {
			panic(fmt.Errorf("pgx level from env: %w", err))
		}
		return l
	}
	return tracelog.LogLevelInfo
}

func slogLevel(level tracelog.LogLevel) slog.Level {
	switch level {
	case tracelog.LogLevelTrace, tracelog.LogLevelDebug:
		return slog.LevelDebug
	case tracelog.LogLevelInfo, tracelog.LogLevelNone:
		return slog.LevelInfo
	case tracelog.LogLevelWarn:
		return slog.LevelWarn
	case tracelog.LogLevelError:
		return slog.LevelError
	default:
		return slog.LevelError
	}
}
