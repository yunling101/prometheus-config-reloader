package log

import (
	"fmt"
	"log/slog"
	"math"
	"os"
	"path/filepath"
	"strings"
)

const dateFormat = "2006-01-02 15:04:05"

func NewLoggerSlog(c Config) (*slog.Logger, error) {
	lvlOption, err := parseLevel(c.Level)
	if err != nil {
		return nil, err
	}

	handler, err := getHandlerFromFormat(c.Format, slog.HandlerOptions{
		Level:       lvlOption,
		AddSource:   true,
		ReplaceAttr: replaceSlogAttributes,
	})
	if err != nil {
		return nil, err
	}

	return slog.New(handler), nil
}

func replaceSlogAttributes(_ []string, a slog.Attr) slog.Attr {
	if a.Key == "time" {
		return slog.Attr{
			Key:   "ts",
			Value: slog.StringValue(a.Value.Time().UTC().Format(dateFormat)),
		}
	}

	if a.Key == "level" {
		return slog.Attr{
			Key:   "level",
			Value: slog.StringValue(strings.ToLower(a.Value.String())),
		}
	}

	if a.Key == "source" {
		return slog.Attr{
			Key:   "caller",
			Value: getCaller(a.Value),
		}
	}

	return a
}

func getCaller(value slog.Value) slog.Value {
	parts := strings.Split(strings.ReplaceAll(value.String(), "}", ""), " ")
	if len(parts) == 3 {
		value = slog.StringValue(fmt.Sprintf("%s:%v", filepath.Base(parts[1]), parts[2]))
	}
	return value
}

func getHandlerFromFormat(format string, opts slog.HandlerOptions) (slog.Handler, error) {
	var handler slog.Handler
	switch strings.ToLower(format) {
	case FormatLogFmt:
		handler = slog.NewTextHandler(os.Stdout, &opts)
		return handler, nil
	case FormatJSON:
		handler = slog.NewJSONHandler(os.Stdout, &opts)
		return handler, nil
	default:
		return nil, fmt.Errorf("log format %s unknown, %v are possible values", format, AvailableLogFormats)
	}
}

func parseLevel(lvl string) (slog.Level, error) {
	switch strings.ToLower(lvl) {
	case LevelAll:
		return slog.LevelDebug, nil
	case LevelDebug:
		return slog.LevelDebug, nil
	case LevelInfo:
		return slog.LevelInfo, nil
	case LevelWarn:
		return slog.LevelWarn, nil
	case LevelError:
		return slog.LevelError, nil
	case LevelNone:
		return math.MaxInt, nil
	default:
		return math.MaxInt, fmt.Errorf("log log_level %s unknown, %v are possible values", lvl, AvailableLogLevels)
	}
}
