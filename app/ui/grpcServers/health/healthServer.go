package server

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	health "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

type helathServer struct{}

func NewHealthGrpcServer(g *grpc.Server) {
	healthServer := helathServer{}
	health.RegisterHealthServer(g, &healthServer)
}

func (s *helathServer) Check(ctx context.Context, req *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
	return &health.HealthCheckResponse{
		Status: health.HealthCheckResponse_SERVING,
	}, nil
}

func (s *helathServer) Watch(in *health.HealthCheckRequest, _ health.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "unimplemented")
}
