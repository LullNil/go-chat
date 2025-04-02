package logger

import (
	"io"
	"log/slog"
	"os"

	"go-chat/internal/lib/logger/handlers/slogpretty"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// GlobalLogger is a global logger for the entire application.
var GlobalLogger *slog.Logger

// SetupLogger initializes the global logger depending on the environment.
func SetupLogger(env string, out io.Writer) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		opts := slogpretty.PrettyHandlerOptions{
			SlogOpts: &slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		}
		handler := opts.NewPrettyHandler(out)
		log = slog.New(handler)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(out, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(out, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log = slog.New(
			slog.NewJSONHandler(out, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	GlobalLogger = log
	return log
}

// GetLogger returns the global logger. If it has not been initialized
// it uses os.Stdout.
func GetLogger() *slog.Logger {
	if GlobalLogger == nil {
		return SetupLogger("prod", os.Stdout)
	}
	return GlobalLogger
}
