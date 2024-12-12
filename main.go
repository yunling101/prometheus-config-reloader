package main

import (
	"context"
	"fmt"
	"github.com/alecthomas/kingpin/v2"
	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/thanos-io/thanos/pkg/reloader"
	logging "github.com/yunling101/prometheus-config-reloader/log"
	"github.com/yunling101/prometheus-config-reloader/metrics"
	"github.com/yunling101/prometheus-config-reloader/version"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	defaultWatchInterval = 3 * time.Minute
)

func main() {
	app := kingpin.New("prometheus-config-reloader", "")
	cfgFile := app.Flag("config-file", "config file watched by the reloader").String()
	watchInterval := app.Flag("watch-interval", "how often the reloader re-reads the configuration file and directories; when set to 0, the program runs only once and exits").Default(defaultWatchInterval.String()).Duration()
	watchedDir := app.Flag("watched-dir", "directory to watch non-recursively").Strings()
	reloadURL := app.Flag("reload-url", "URL to trigger the configuration").Default("http://127.0.0.1:9090/-/reload").URL()
	listenAddress := app.Flag(
		"listen-address",
		"address on which to expose metrics (disabled when empty)").Default("0.0.0.0:9096").String()
	var logConfig logging.Config
	app.Flag(
		"log-format",
		fmt.Sprintf("log format to use. Possible values: %s", strings.Join(logging.AvailableLogFormats, ", "))).
		Default(logging.FormatLogFmt).StringVar(&logConfig.Format)
	app.Flag(
		"log-level",
		fmt.Sprintf("log level to use. Possible values: %s", strings.Join(logging.AvailableLogLevels, ", "))).
		Default(logging.LevelInfo).StringVar(&logConfig.Level)

	if _, err := app.Parse(os.Args[1:]); err != nil {
		_, _ = fmt.Fprintln(os.Stdout, err)
		os.Exit(2)
	}

	logger, err := logging.NewLoggerSlog(logConfig)
	if err != nil {
		log.Fatal(err)
	}

	logger.Info("Starting prometheus-config-reloader", "version", version.Info())

	r := metrics.NewRegistry("prometheus_config_reloader")
	var (
		g           run.Group
		ctx, cancel = context.WithCancel(context.Background())
	)

	{
		kitLogger, err := logging.NewLogger(logConfig)
		if err != nil {
			log.Fatal(err)
		}

		rel := reloader.New(kitLogger, r, &reloader.Options{
			CfgFile:                       *cfgFile,
			WatchedDirs:                   *watchedDir,
			WatchInterval:                 *watchInterval,
			ReloadURL:                     *reloadURL,
			RetryInterval:                 defaultWatchInterval,
			TolerateEnvVarExpansionErrors: true,
		})
		g.Add(func() error {
			return rel.Watch(ctx)
		}, func(error) {
			cancel()
		})
	}

	{
		http.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{Registry: r}))
		http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":"up"}`))
		})

		srv := &http.Server{Addr: *listenAddress}

		g.Add(func() error {
			logger.Info("Starting web server for metrics", "listen", *listenAddress)
			return srv.ListenAndServe()
		}, func(error) {
			_ = srv.Close()
		})
	}

	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	g.Add(func() error {
		select {
		case <-term:
			logger.Info("Received SIGTERM, exiting gracefully...")
		case <-ctx.Done():
		}
		return nil
	}, func(error) {})

	if err = g.Run(); err != nil {
		logger.Error("Failed to run", "err", err)
		os.Exit(1)
	}
}
