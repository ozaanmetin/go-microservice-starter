package circuitbreaker

import (
	"time"

	"github.com/sony/gobreaker"
)

// CircuitBreaker wraps gobreaker.CircuitBreaker to provide a clean interface
type CircuitBreaker struct {
	cb *gobreaker.CircuitBreaker
}

// Config holds circuit breaker configuration
type Config struct {
	Name          string
	MaxRequests   uint32
	Interval      time.Duration
	Timeout       time.Duration
	ReadyToTrip   func(counts gobreaker.Counts) bool
	OnStateChange func(name string, from gobreaker.State, to gobreaker.State)
}

// New creates a new CircuitBreaker instance
func NewCircuitBreaker(cfg Config) *CircuitBreaker {
	settings := gobreaker.Settings{
		Name:          cfg.Name,
		MaxRequests:   cfg.MaxRequests,
		Interval:      cfg.Interval,
		Timeout:       cfg.Timeout,
		ReadyToTrip:   cfg.ReadyToTrip,
		OnStateChange: cfg.OnStateChange,
	}

	return &CircuitBreaker{
		cb: gobreaker.NewCircuitBreaker(settings),
	}
}

// Execute runs the given function with circuit breaker protection
func (cb *CircuitBreaker) Execute(fn func() (interface{}, error)) (interface{}, error) {
	return cb.cb.Execute(fn)
}

// State returns the current state of the circuit breaker
func (cb *CircuitBreaker) State() gobreaker.State {
	return cb.cb.State()
}

// Name returns the name of the circuit breaker
func (cb *CircuitBreaker) Name() string {
	return cb.cb.Name()
}
