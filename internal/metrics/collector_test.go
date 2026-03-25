package metrics_test

import (
	"testing"

	"github.com/Chalupa-Tech/go-telemetry-mesh/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

func TestNewCollector(t *testing.T) {
	reg := prometheus.NewRegistry()
	c := metrics.NewCollector(reg)

	if c.ProbeDuration == nil {
		t.Fatal("ProbeDuration histogram should not be nil")
	}
	if c.ProbeTotal == nil {
		t.Fatal("ProbeTotal counter should not be nil")
	}
	if c.ProbeErrors == nil {
		t.Fatal("ProbeErrors counter should not be nil")
	}

	// Verify metrics are registered by gathering
	families, err := reg.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	names := make(map[string]bool)
	for _, f := range families {
		names[f.GetName()] = true
	}

	for _, want := range []string{
		"telemetry_mesh_peers_discovered",
		"telemetry_mesh_cluster_healthy",
	} {
		if !names[want] {
			t.Errorf("Expected metric %q to be registered", want)
		}
	}
}

func TestRecordSuccess(t *testing.T) {
	reg := prometheus.NewRegistry()
	c := metrics.NewCollector(reg)

	c.RecordSuccess("node-a", "node-b", 0.001)
	c.RecordSuccess("node-a", "node-b", 0.002)

	families, err := reg.Gather()
	if err != nil {
		t.Fatalf("Failed to gather: %v", err)
	}

	var found bool
	for _, f := range families {
		if f.GetName() == "telemetry_mesh_probe_total" {
			for _, m := range f.GetMetric() {
				for _, l := range m.GetLabel() {
					if l.GetName() == "result" && l.GetValue() == "success" {
						if m.GetCounter().GetValue() != 2 {
							t.Errorf("Expected 2 success probes, got %v", m.GetCounter().GetValue())
						}
						found = true
					}
				}
			}
		}
	}
	if !found {
		t.Error("Success counter metric not found")
	}
}

func TestRecordError(t *testing.T) {
	reg := prometheus.NewRegistry()
	c := metrics.NewCollector(reg)

	c.RecordError("node-a", "node-b", "timeout")

	families, err := reg.Gather()
	if err != nil {
		t.Fatalf("Failed to gather: %v", err)
	}

	var found bool
	for _, f := range families {
		if f.GetName() == "telemetry_mesh_probe_errors_total" {
			for _, m := range f.GetMetric() {
				for _, l := range m.GetLabel() {
					if l.GetName() == "error_type" && l.GetValue() == "timeout" {
						if m.GetCounter().GetValue() != 1 {
							t.Errorf("Expected 1 error, got %v", m.GetCounter().GetValue())
						}
						found = true
					}
				}
			}
		}
	}
	if !found {
		t.Error("Error counter metric not found")
	}
}
