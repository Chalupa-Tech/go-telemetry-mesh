package server

import (
	"context"
	"log/slog"
	"time"

	meshv1 "github.com/Chalupa-Tech/go-telemetry-mesh/api/proto/v1"
)

type MeshServer struct {
	meshv1.UnimplementedMeshServiceServer
}

func NewMeshServer() *MeshServer {
	return &MeshServer{}
}

func (s *MeshServer) Ping(ctx context.Context, req *meshv1.PingRequest) (*meshv1.PingResponse, error) {
	// Log the ping
	slog.Info("Received Ping",
		"origin", req.Origin,
		"client_timestamp", req.Timestamp,
	)

	return &meshv1.PingResponse{
		Success:    true,
		ReceivedAt: time.Now().UnixNano(),
	}, nil
}
