package log

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-kit/log"
	loglevel "github.com/go-kit/log/level"
)

const (
	LevelAll   = "all"
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
	LevelNone  = "none"
)

const (
	FormatLogFmt = "logfmt"
	FormatJSON   = "json"
)

type Config struct {
	Level  string
	Format string
}

var (
	defaultDateTime = log.TimestampFormat(
		func() time.Time { return time.Now().UTC() },
		dateFormat,
	)
)

func RegisterFlags(fs *flag.FlagSet, c *Config) {
	fs.StringVar(&c.Level, "log-level", "info",
		fmt.Sprintf("Log level to use. Possible values: %s",
			strings.Join(AvailableLogLevels, ", ")))
	fs.StringVar(&c.Format, "log-format", "logfmt",
		fmt.Sprintf("Log format to use. Possible values: %s",
			strings.Join(AvailableLogFormats, ", ")))
}

func NewLogger(c Config) (log.Logger, error) {
	var (
		logger    log.Logger
		lvlOption loglevel.Option
	)

	switch strings.ToLower(c.Level) {
	case LevelAll:
		lvlOption = loglevel.AllowAll()
	case LevelDebug:
		lvlOption = loglevel.AllowDebug()
	case LevelInfo:
		lvlOption = loglevel.AllowInfo()
	case LevelWarn:
		lvlOption = loglevel.AllowWarn()
	case LevelError:
		lvlOption = loglevel.AllowError()
	case LevelNone:
		lvlOption = loglevel.AllowNone()
	default:
		return nil, fmt.Errorf("log log_level %s unknown, %v are possible values", c.Level, AvailableLogLevels)
	}

	switch c.Format {
	case FormatLogFmt:
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	case FormatJSON:
		logger = log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	default:
		return nil, fmt.Errorf("log format %s unknown, %v are possible values", c.Format, AvailableLogFormats)
	}

	logger = log.With(logger, "ts", defaultDateTime)
	logger = loglevel.NewFilter(logger, lvlOption)
	logger = log.With(logger, "caller", log.DefaultCaller)

	return logger, nil
}

var AvailableLogLevels = []string{
	LevelAll,
	LevelDebug,
	LevelInfo,
	LevelWarn,
	LevelError,
	LevelNone,
}

var AvailableLogFormats = []string{
	FormatLogFmt,
	FormatJSON,
}
