package exporter

import "github.com/prometheus/client_golang/prometheus"

type collector struct {
	gearman       *gearman
	up            *prometheus.Desc
	statusTotal   *prometheus.Desc
	statusRunning *prometheus.Desc
	statusWorkers *prometheus.Desc
}

// based on https://github.com/hnlq715/nginx-vts-exporter/
func newFuncMetric(metricName string, docString string, labels []string) *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName(metricsNamespace, "", metricName),
		docString, labels, nil,
	)
}

func newCollector(g *gearman) *collector {
	return &collector{
		gearman:       g,
		up:            newFuncMetric("up", "is gearman up", []string{"version"}),
		statusTotal:   newFuncMetric("status_total", "number of jobs in the queue", []string{"function"}),
		statusRunning: newFuncMetric("status_running", "number of running jobs", []string{"function"}),
		statusWorkers: newFuncMetric("status_workers", "number of number of capable workers", []string{"function"}),
	}
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up
	ch <- c.statusTotal
	ch <- c.statusRunning
	ch <- c.statusWorkers
}

func (c *collector) collectVersion(ch chan<- prometheus.Metric) {
	up := 1.0
	// TODO: should we log this error?
	v, err := c.gearman.getVersion()
	if err != nil {
		up = 0.0
		v = "unknown"
	}
	ch <- prometheus.MustNewConstMetric(
		c.up,
		prometheus.GaugeValue,
		up,
		v)
}
func (c *collector) collectStatus(ch chan<- prometheus.Metric) {
	s, err := c.gearman.getStatus()
	if err != nil {
		// TODO: should we log this error?
		return
	}

	for k, v := range s {
		ch <- prometheus.MustNewConstMetric(
			c.statusTotal,
			prometheus.GaugeValue,
			float64(v.total),
			k)

		ch <- prometheus.MustNewConstMetric(
			c.statusRunning,
			prometheus.GaugeValue,
			float64(v.running),
			k)

		ch <- prometheus.MustNewConstMetric(
			c.statusWorkers,
			prometheus.GaugeValue,
			float64(v.workers),
			k)
	}
}

func (c *collector) Collect(ch chan<- prometheus.Metric) {
	c.collectVersion(ch)
	c.collectStatus(ch)
}
