package services

import (
	"github.com/prometheus/client_golang/prometheus"
)

type MetricsService struct {
	HttpTotalRequests        *prometheus.CounterVec
	HttpSuccessfulRequests   *prometheus.CounterVec
	HttpUnsuccessfulRequests *prometheus.CounterVec
	HttpRequestDuration      *prometheus.HistogramVec
	AverageRequestDuration   *prometheus.GaugeVec
	RequestsPerTimeUnit      *prometheus.CounterVec
	Registry                 *prometheus.Registry
}

func NewMetricsService() *MetricsService {
	registry := prometheus.NewRegistry()

	httpTotalRequests := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_total_requests",
			Help: "Total number of HTTP requests in last 24h",
		},
		[]string{},
	)
	registry.MustRegister(httpTotalRequests)

	httpSuccessfulRequests := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_successful_requests",
			Help: "Number of successful HTTP requests in last 24h (2xx, 3xx).",
		},
		[]string{},
	)
	registry.MustRegister(httpSuccessfulRequests)

	httpUnsuccessfulRequests := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_unsuccessful_requests",
			Help: "Number of unsuccessful HTTP requests in last 24h (4xx, 5xx).",
		},
		[]string{},
	)
	registry.MustRegister(httpUnsuccessfulRequests)

	averageRequestDuration := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "average_request_duration_seconds",
			Help: "Average request duration for each endpoint.",
		},
		[]string{"method", "endpoint"},
	)
	registry.MustRegister(averageRequestDuration)

	requestsPerTimeUnit := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests_per_time_unit",
			Help: "Number of requests per time unit (e.g., per minute or per second) for each endpoint.",
		},
		[]string{"method", "endpoint", "time_unit"},
	)
	registry.MustRegister(requestsPerTimeUnit)

	return &MetricsService{
		HttpTotalRequests:        httpTotalRequests,
		HttpSuccessfulRequests:   httpSuccessfulRequests,
		HttpUnsuccessfulRequests: httpUnsuccessfulRequests,
		AverageRequestDuration:   averageRequestDuration,
		RequestsPerTimeUnit:      requestsPerTimeUnit,
		Registry:                 registry,
	}
}
