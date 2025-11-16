package healthcheck

import (
	"context"
	"time"
)

// HealthCheckRequest represents the health check request (empty)
type HealthCheckRequest struct{}

// HealthCheckResponse represents the health check response
type HealthCheckResponse struct {
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`
}

// NewHealthCheckHandler creates a new health check handler
func NewHealthCheckHandler() *HealthCheckHandler {
	return &HealthCheckHandler{}
}


// HealthCheckHandler handles health check requests
type HealthCheckHandler struct{}

// Handle implements the HandlerInterface for health checks
func (h *HealthCheckHandler) Handle(ctx context.Context, req *HealthCheckRequest) (*HealthCheckResponse, error) {
	return &HealthCheckResponse{
		Status:    "Ok",
		Timestamp: time.Now().Unix(),
	}, nil
}
