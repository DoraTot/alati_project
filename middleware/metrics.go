package middleware

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"projekat/services"
	"time"
)

type Metrics struct {
	service *services.MetricsService
}

func NewMetrics(service *services.MetricsService) *Metrics {
	return &Metrics{service}
}

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (r *ResponseWriter) WriteHeader(status int) {
	r.statusCode = status
	r.ResponseWriter.WriteHeader(status)
}

func (m *Metrics) Count(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Scrapping...")
		start := time.Now()

		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		rw := &ResponseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()

		statusCode := rw.statusCode
		if statusCode >= 200 && statusCode < 400 {
			m.service.HttpSuccessfulRequests.WithLabelValues().Inc()
		} else if statusCode >= 400 && statusCode < 600 {
			m.service.HttpUnsuccessfulRequests.WithLabelValues().Inc()
		}

		m.service.HttpTotalRequests.WithLabelValues().Inc()
		m.service.AverageRequestDuration.WithLabelValues(r.Method, path).Set(duration)
		m.service.RequestsPerTimeUnit.WithLabelValues(r.Method, path, "seconds").Inc()
	})
}

func (m *Metrics) MetricsHandler() http.Handler {
	return promhttp.HandlerFor(m.service.Registry, promhttp.HandlerOpts{})
}

func AdaptPrometheusHandler(handler http.Handler, metrics *Metrics) http.Handler {
	return metrics.Count(handler)
}
