package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	meshv1 "github.com/Chalupa-Tech/go-telemetry-mesh/api/proto/v1"
	"github.com/Chalupa-Tech/go-telemetry-mesh/internal/client"
	"github.com/Chalupa-Tech/go-telemetry-mesh/internal/metrics"
	"github.com/Chalupa-Tech/go-telemetry-mesh/internal/server"
	telemetry "github.com/Chalupa-Tech/go-telemetry"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	version := os.Getenv("APP_VERSION")
	if version == "" {
		version = "dev"
	}
	shutdown, err := telemetry.Init(ctx, "go-telemetry-mesh", version)
	if err != nil {
		slog.Error("failed to init telemetry", "error", err)
		// Non-fatal: continue without tracing (gRPC otelgrpc still works via global provider).
	} else {
		defer shutdown()
	}

	cfg := loadConfig()
	slog.Info("Starting telemetry-mesh",
		"node", cfg.nodeName,
		"grpc_port", cfg.grpcPort,
		"http_port", cfg.httpPort,
		"interval", cfg.pingInterval,
	)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	// Metrics
	reg := prometheus.DefaultRegisterer
	collector := metrics.NewCollector(reg)

	// gRPC server
	grpcLis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.grpcPort))
	if err != nil {
		slog.Error("Failed to listen for gRPC", "error", err)
		os.Exit(1)
	}
	grpcSrv := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))
	meshv1.RegisterMeshServiceServer(grpcSrv, server.NewMeshServer(cfg.nodeName))

	// HTTP server (metrics + health)
	ready := &atomic.Bool{}
	httpSrv := &http.Server{
		Addr:              fmt.Sprintf(":%s", cfg.httpPort),
		Handler:           server.NewHTTPMux(ready),
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Mesh client
	meshClient := client.NewMeshClient(
		cfg.headlessSvc, cfg.grpcPort, cfg.nodeName, cfg.podIP,
		cfg.pingInterval, cfg.pingTimeout, collector,
	)

	// Start gRPC server
	go func() {
		slog.Info("gRPC server listening", "addr", grpcLis.Addr())
		if err := grpcSrv.Serve(grpcLis); err != nil {
			slog.Error("gRPC server error", "error", err)
		}
	}()

	// Start HTTP server
	go func() {
		slog.Info("HTTP server listening", "addr", httpSrv.Addr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP server error", "error", err)
		}
	}()

	// Mark ready and start probing
	ready.Store(true)
	go meshClient.Start(ctx)

	<-ctx.Done()
	slog.Info("Shutting down...")

	// Graceful shutdown
	grpcSrv.GracefulStop()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	httpSrv.Shutdown(shutdownCtx)

	slog.Info("Shutdown complete")
}

type config struct {
	nodeName     string
	podIP        string
	grpcPort     string
	httpPort     string
	headlessSvc  string
	pingInterval time.Duration
	pingTimeout  time.Duration
}

func loadConfig() config {
	return config{
		nodeName:     envOrDefault("NODE_NAME", "unknown"),
		podIP:        envOrDefault("POD_IP", ""),
		grpcPort:     envOrDefault("GRPC_PORT", "9090"),
		httpPort:     envOrDefault("HTTP_PORT", "8080"),
		headlessSvc:  envOrDefault("HEADLESS_SVC", "telemetry-mesh-headless"),
		pingInterval: durationOrDefault("PING_INTERVAL", 10*time.Second),
		pingTimeout:  durationOrDefault("PING_TIMEOUT", 5*time.Second),
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func durationOrDefault(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			slog.Warn("Invalid duration, using default", "key", key, "value", v, "default", fallback)
			return fallback
		}
		return d
	}
	return fallback
}

func logLevel() slog.Level {
	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
