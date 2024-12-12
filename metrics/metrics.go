package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/collectors/version"
	"regexp"
)

func NewRegistry(program string) *prometheus.Registry {
	r := prometheus.NewRegistry()

	r.MustRegister(
		collectors.NewGoCollector(
			collectors.WithGoCollectorRuntimeMetrics(
				collectors.MetricsGC,
				collectors.MetricsScheduler,
				collectors.MetricsMemory,
				collectors.GoRuntimeMetricsRule{Matcher: regexp.MustCompile(`^/sync/.*`)},
			),
		),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		version.NewCollector(program),
	)

	return r
}
