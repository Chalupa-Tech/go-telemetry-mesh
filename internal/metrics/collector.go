package metrics

import "github.com/prometheus/client_golang/prometheus"

// Collector holds all Prometheus metrics for the telemetry mesh.
type Collector struct {
	ProbeDuration    *prometheus.HistogramVec
	ProbeTotal       *prometheus.CounterVec
	ProbeErrors      *prometheus.CounterVec
	PeersDiscovered  prometheus.Gauge
	ClusterHealthy   prometheus.Gauge
}

// NewCollector creates and registers all mesh metrics.
func NewCollector(reg prometheus.Registerer) *Collector {
	c := &Collector{
		ProbeDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "telemetry_mesh_probe_duration_seconds",
			Help:    "Latency of mesh probe pings between nodes.",
			Buckets: []float64{0.0005, 0.001, 0.0025, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0},
		}, []string{"source_node", "target_node"}),

		ProbeTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "telemetry_mesh_probe_total",
			Help: "Total number of mesh probes sent.",
		}, []string{"source_node", "target_node", "result"}),

		ProbeErrors: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "telemetry_mesh_probe_errors_total",
			Help: "Total number of mesh probe errors by type.",
		}, []string{"source_node", "target_node", "error_type"}),

		PeersDiscovered: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "telemetry_mesh_peers_discovered",
			Help: "Number of peers currently discovered via DNS.",
		}),

		ClusterHealthy: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "telemetry_mesh_cluster_healthy",
			Help: "1 if all peers are reachable, 0 otherwise.",
		}),
	}

	reg.MustRegister(c.ProbeDuration, c.ProbeTotal, c.ProbeErrors, c.PeersDiscovered, c.ClusterHealthy)
	return c
}

// RecordSuccess records a successful probe with its latency.
func (c *Collector) RecordSuccess(source, target string, durationSec float64) {
	c.ProbeDuration.WithLabelValues(source, target).Observe(durationSec)
	c.ProbeTotal.WithLabelValues(source, target, "success").Inc()
}

// RecordError records a failed probe with the error type.
func (c *Collector) RecordError(source, target, errType string) {
	c.ProbeTotal.WithLabelValues(source, target, "error").Inc()
	c.ProbeErrors.WithLabelValues(source, target, errType).Inc()
}
