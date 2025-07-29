package log

import (
	"log/slog"
	"os"
)

var (
	logLevel      = new(slog.LevelVar)
	currentFormat = "text"
)

func init() {
	level, _ := parseLogLevel(os.Getenv("LOG_LEVEL"))
	logLevel.Set(level)

	updateLogger(os.Getenv("LOG_FORMAT"))
}
