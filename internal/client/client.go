package client

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"

	meshv1 "github.com/Chalupa-Tech/go-telemetry-mesh/api/proto/v1"
	"github.com/Chalupa-Tech/go-telemetry-mesh/internal/metrics"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// MeshClient discovers peers via DNS and pings them over gRPC.
type MeshClient struct {
	headlessSvc string
	port        string
	nodeName    string
	interval    time.Duration
	timeout     time.Duration
	metrics     *metrics.Collector
	podIP       string
}

// NewMeshClient creates a new mesh client.
func NewMeshClient(headlessSvc, port, nodeName, podIP string, interval, timeout time.Duration, m *metrics.Collector) *MeshClient {
	return &MeshClient{
		headlessSvc: headlessSvc,
		port:        port,
		nodeName:    nodeName,
		podIP:       podIP,
		interval:    interval,
		timeout:     timeout,
		metrics:     m,
	}
}

// Start runs the ping loop until the context is cancelled.
func (c *MeshClient) Start(ctx context.Context) {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	c.scanAndPing(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.scanAndPing(ctx)
		}
	}
}

func (c *MeshClient) scanAndPing(ctx context.Context) {
	ips, err := net.LookupHost(c.headlessSvc)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to lookup peers", "error", err, "service", c.headlessSvc)
		return
	}

	// Filter out self
	var peers []string
	for _, ip := range ips {
		if ip != c.podIP {
			peers = append(peers, ip)
		}
	}

	c.metrics.PeersDiscovered.Set(float64(len(peers)))

	var (
		wg       sync.WaitGroup
		mu       sync.Mutex
		failures int
	)

	for _, ip := range peers {
		wg.Add(1)
		go func(targetIP string) {
			defer wg.Done()
			if !c.pingPeer(ctx, targetIP) {
				mu.Lock()
				failures++
				mu.Unlock()
			}
		}(ip)
	}
	wg.Wait()

	if failures == 0 && len(peers) > 0 {
		c.metrics.ClusterHealthy.Set(1)
	} else {
		c.metrics.ClusterHealthy.Set(0)
	}
}

// pingPeer sends a gRPC ping and records metrics. Returns true on success.
func (c *MeshClient) pingPeer(ctx context.Context, ip string) bool {
	target := fmt.Sprintf("%s:%s", ip, c.port)

	conn, err := grpc.NewClient(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to dial peer", "peer", ip, "error", err)
		c.metrics.RecordError(c.nodeName, ip, "connection")
		return false
	}
	defer conn.Close()

	client := meshv1.NewMeshServiceClient(conn)

	callCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	req := &meshv1.PingRequest{
		Origin:    c.nodeName,
		Timestamp: time.Now().UnixNano(),
	}

	start := time.Now()
	resp, err := client.Ping(callCtx, req)
	duration := time.Since(start)

	if err != nil {
		errType := "rpc"
		if callCtx.Err() == context.DeadlineExceeded {
			errType = "timeout"
		}
		slog.WarnContext(ctx, "Ping failed", "peer", ip, "error", err, "duration", duration)
		c.metrics.RecordError(c.nodeName, ip, errType)
		return false
	}

	targetName := resp.NodeName
	if targetName == "" {
		targetName = ip
	}
	c.metrics.RecordSuccess(c.nodeName, targetName, duration.Seconds())
	slog.DebugContext(ctx, "Ping success",
		"peer", targetName,
		"duration", duration,
		"remote_ts", resp.ReceivedAt,
	)
	return true
}
