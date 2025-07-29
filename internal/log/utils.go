package log

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

func parseLogLevel(levelStr string) (slog.Level, error) {
	switch strings.ToUpper(levelStr) {
	case "DEBUG":
		return slog.LevelDebug, nil
	case "INFO":
		return slog.LevelInfo, nil
	case "WARN":
		return slog.LevelWarn, nil
	case "ERROR":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, fmt.Errorf("invalid log level: %s. valid levels are: DEBUG, INFO, WARN, ERROR", levelStr)
	}
}

func updateLogger(format string) {
	if format == "" {
		format = currentFormat
	}

	var handler slog.Handler
	handlerOpts := &slog.HandlerOptions{Level: logLevel}

	if strings.ToLower(format) == "json" {
		handler = slog.NewJSONHandler(os.Stdout, handlerOpts)
		currentFormat = "json"
	} else {
		handler = slog.NewTextHandler(os.Stdout, handlerOpts)
		currentFormat = "text"
	}
	slog.SetDefault(slog.New(handler))
}

// SetLevel sets the logging level dynamically.
func SetLevel(levelStr string) error {
	level, err := parseLogLevel(levelStr)
	if err != nil {
		return err
	}
	logLevel.Set(level)
	return nil
}

// SetFormat sets the logging format dynamically.
func SetFormat(format string) error {
	lowerFormat := strings.ToLower(format)
	if lowerFormat != "json" && lowerFormat != "text" {
		return fmt.Errorf("invalid log format: %s. valid formats are: json, text", format)
	}
	updateLogger(lowerFormat)
	return nil
}

func GetLevel() string {
	return logLevel.Level().String()
}

func GetFormat() string {
	return currentFormat
}
