package circuitBreakerExample

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/ozaanmetin/go-microservice-starter/pkg/circuitbreaker"
	appErrors "github.com/ozaanmetin/go-microservice-starter/pkg/errors"
	"github.com/ozaanmetin/go-microservice-starter/pkg/logging"
	"github.com/sony/gobreaker"
)

// ExampleRequest represents the example request
type ExampleRequest struct{}

// ExampleResponse represents the example response
type ExampleResponse struct {
	Message        string `json:"message"`
	CircuitBreaker string `json:"circuit_breaker"`
}

// ExampleHandler demonstrates circuit breaker usage
type ExampleHandler struct {
	cb *circuitbreaker.CircuitBreaker
}

// NewExampleHandler creates a new example handler with circuit breaker
func NewExampleHandler() *ExampleHandler {
	cbSettings := circuitbreaker.Config{
			Name: "example-circuit-breaker",
			MaxRequests:   3,
			Interval:      5 * time.Second,
			Timeout:       10 * time.Second,
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
				return counts.Requests >= 3 && failureRatio >= 0.6
			},
			OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
				logging.L().
				WithField("name", name).
				WithField("from", from.String()).
				WithField("to", to.String()).
				Info("Circuit breaker state changed")
			},
		}

	return &ExampleHandler{
		cb: circuitbreaker.NewCircuitBreaker(cbSettings),
	}
}

// Handle implements the HandlerInterface for example requests
func (h *ExampleHandler) Handle(ctx context.Context, req *ExampleRequest) (*ExampleResponse, error) {
	// Execute operation with circuit breaker protection
	result, err := h.cb.Execute(func() (interface{}, error) {
		return h.callExternalService()
	})

	if err != nil {
		// Circuit breaker error
		return nil, appErrors.NewServiceUnavailableError("Service temporarily unavailable", err)
	}

	message := result.(string)
	return &ExampleResponse{
		Message:        message,
		CircuitBreaker: h.cb.State().String(),
	}, nil
}

// callExternalService simulates an external service call that may fail
func (h *ExampleHandler) callExternalService() (string, error) {
	// Simulate random failures for demonstration
	if rand.Intn(10) < 5 {
		return "", errors.New("simulated external service failure")
	}
	return "External service call successful", nil
}
