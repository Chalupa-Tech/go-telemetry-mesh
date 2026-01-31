package client

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"sync"
	"time"

	meshv1 "github.com/Chalupa-Tech/go-telemetry-mesh/api/proto/v1"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MeshClient struct {
	HeadlessSvc string
	Port        string
	MyOrigin    string
}

func NewMeshClient(headlessSvc, port string) *MeshClient {
	hostname, _ := os.Hostname()
	return &MeshClient{
		HeadlessSvc: headlessSvc,
		Port:        port,
		MyOrigin:    hostname,
	}
}

func (c *MeshClient) Start(ctx context.Context) {
	ticker := time.NewTicker(15 * time.Second) // Ping every 15s
	defer ticker.Stop()

	// Initial run
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
	ips, err := net.LookupHost(c.HeadlessSvc)
	if err != nil {
		slog.Error("Failed to lookup peers", "error", err, "service", c.HeadlessSvc)
		return
	}

	var wg sync.WaitGroup
	for _, ip := range ips {
		// Skip self? Hard to detect IP matching reliable without more logic,
		// but pinging self is fine for now.
		wg.Add(1)
		go func(targetIP string) {
			defer wg.Done()
			c.pingPeer(ctx, targetIP)
		}(ip)
	}
	wg.Wait()
}

func (c *MeshClient) pingPeer(ctx context.Context, ip string) {
	// Create connection with OTel instrumentation
	target := fmt.Sprintf("%s:%s", ip, c.Port)

	// Note: Creating a new client for every ping is inefficient for high scale,
	// but acceptable for a canary mesh with standard peer discovery intervals.
	// Ideally we would cache connections.
	conn, err := grpc.NewClient(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		slog.Error("Failed to dial peer", "peer", ip, "error", err)
		return
	}
	defer conn.Close()

	client := meshv1.NewMeshServiceClient(conn)

	req := &meshv1.PingRequest{
		Origin:    c.MyOrigin,
		Timestamp: time.Now().UnixNano(),
	}

	// Context with timeout for the call
	callCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := client.Ping(callCtx, req)
	if err != nil {
		slog.Error("Ping failed", "peer", ip, "error", err)
		return
	}

	slog.Info("Ping success", "peer", ip, "latency_ns", time.Now().UnixNano()-req.Timestamp, "remote_ts", resp.ReceivedAt)
}
