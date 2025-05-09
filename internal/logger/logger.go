package logger

import (
	"log/slog"
	"os"
	"time"

	"github.com/aik27/otus_go_image_previewer/internal/config"
)

func SetupLogger(cnf *config.Config) {
	var logLevel slog.Level

	switch cnf.App.LogLevel {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		// @TODO output message in log format (we dont have struct logger at this point)
		panic("invalid log level")
	}

	stdout := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:       logLevel,
		ReplaceAttr: replaceAttr,
	}))

	slog.SetDefault(stdout)
}

func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	_ = groups
	switch a.Key {
	case slog.TimeKey:
		return slog.String("timestamp", a.Value.Time().Format(time.RFC3339))
	case slog.MessageKey:
		return slog.String("rest", a.Value.String())
	case slog.LevelKey:
		return slog.String("severity", a.Value.String())
	default:
		return a
	}
}
