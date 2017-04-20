package exporter

import "github.com/prometheus/client_golang/prometheus"

const metricsNamespace = "gearman"

var (
	requestsHistogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: metricsNamespace,
			Name:      "requests_duration",
			Help:      "duration of gearman metrics requests in seconds",
		},
	)

	gearmanUp = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: metricsNamespace,
			Name:      "up",
			Help:      "if gearman is up",
		},
		[]string{"version"},
	)
)

func init() {
	//prometheus.MustRegister(requestsHistogram)
	//prometheus.MustRegister(gearmanUp)
}
