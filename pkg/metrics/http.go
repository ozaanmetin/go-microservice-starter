package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// httpRequestsTotal counts HTTP requests by method, path, and status
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// httpRequestDuration tracks HTTP request duration in seconds
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: []float64{0.001, 0.01, 0.1, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"method", "path", "status"},
	)

	// httpRequestsInFlight tracks current number of requests being processed
	httpRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Current number of HTTP requests being processed",
		},
	)
)

// RecordHTTPRequest records metrics for an HTTP request
func RecordHTTPRequest(method, path string, statusCode int, duration time.Duration) {
	status := strconv.Itoa(statusCode)
	httpRequestsTotal.WithLabelValues(method, path, status).Inc()
	httpRequestDuration.WithLabelValues(method, path, status).Observe(duration.Seconds())
}

// IncrementInFlight increments the in-flight requests counter
func IncrementInFlight() {
	httpRequestsInFlight.Inc()
}

// DecrementInFlight decrements the in-flight requests counter
func DecrementInFlight() {
	httpRequestsInFlight.Dec()
}
