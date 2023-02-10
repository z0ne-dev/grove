package util

import (
	"cdr.dev/slog"
	"context"
	"github.com/jackc/pgx/v4"
)

func SlogPgxLogger(logger slog.Logger) func(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	return func(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
		log := logger.Debug
		switch level {
		case pgx.LogLevelNone:
			return

		case pgx.LogLevelError:
			log = logger.Error
			break

		case pgx.LogLevelWarn:
			log = logger.Warn
			break

		case pgx.LogLevelInfo:
			log = logger.Info
		}

		fields := make([]slog.Field, 0, len(data))
		for k, v := range data {
			fields = append(fields, slog.F(k, v))
		}

		log(ctx, msg, fields...)
	}
}
