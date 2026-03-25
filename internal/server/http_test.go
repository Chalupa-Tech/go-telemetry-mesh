package server_test

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/Chalupa-Tech/go-telemetry-mesh/internal/server"
)

func TestHealthz(t *testing.T) {
	ready := &atomic.Bool{}
	mux := server.NewHTTPMux(ready)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

func TestReadyzNotReady(t *testing.T) {
	ready := &atomic.Bool{}
	mux := server.NewHTTPMux(ready)

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected 503 when not ready, got %d", w.Code)
	}
}

func TestReadyzReady(t *testing.T) {
	ready := &atomic.Bool{}
	ready.Store(true)
	mux := server.NewHTTPMux(ready)

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 when ready, got %d", w.Code)
	}
}

func TestMetricsEndpoint(t *testing.T) {
	ready := &atomic.Bool{}
	mux := server.NewHTTPMux(ready)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 for /metrics, got %d", w.Code)
	}
	if w.Header().Get("Content-Type") == "" {
		t.Error("Expected Content-Type header on /metrics")
	}
}
