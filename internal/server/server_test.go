package server_test

import (
	"context"
	"testing"

	meshv1 "github.com/Chalupa-Tech/go-telemetry-mesh/api/proto/v1"
	"github.com/Chalupa-Tech/go-telemetry-mesh/internal/server"
)

func TestPing(t *testing.T) {
	srv := server.NewMeshServer("test-node")

	resp, err := srv.Ping(context.Background(), &meshv1.PingRequest{
		Origin:    "other-node",
		Timestamp: 1234567890,
	})
	if err != nil {
		t.Fatalf("Ping returned error: %v", err)
	}
	if !resp.Success {
		t.Error("Expected Success=true")
	}
	if resp.ReceivedAt == 0 {
		t.Error("Expected non-zero ReceivedAt timestamp")
	}
}
