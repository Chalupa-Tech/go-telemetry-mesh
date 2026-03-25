package server

import (
	"context"
	"log/slog"
	"time"

	meshv1 "github.com/Chalupa-Tech/go-telemetry-mesh/api/proto/v1"
)

// MeshServer implements the gRPC MeshService.
type MeshServer struct {
	meshv1.UnimplementedMeshServiceServer
	nodeName string
}

// NewMeshServer creates a new MeshServer identified by nodeName.
func NewMeshServer(nodeName string) *MeshServer {
	return &MeshServer{nodeName: nodeName}
}

func (s *MeshServer) Ping(ctx context.Context, req *meshv1.PingRequest) (*meshv1.PingResponse, error) {
	slog.Debug("Received Ping",
		"origin", req.Origin,
		"node", s.nodeName,
	)

	return &meshv1.PingResponse{
		Success:    true,
		ReceivedAt: time.Now().UnixNano(),
	}, nil
}
