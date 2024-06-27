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
	log.Printf("Setting response status code to: %d", status)
	r.statusCode = status
	r.ResponseWriter.WriteHeader(status)
}

func (m *Metrics) Count(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request received for:", r.URL.Path)
		start := time.Now()

		route := mux.CurrentRoute(r)
		if route == nil {
			log.Println("No current route found for request")
		}

		path, err := route.GetPathTemplate()
		if err != nil {
			log.Printf("Error getting path template: %v", err)
		}
		log.Printf("Path template: %s", path)

		rw := &ResponseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()
		log.Printf("Processed request for %s in %f seconds", path, duration)

		statusCode := rw.statusCode
		log.Printf("Response status code: %d", statusCode)

		if statusCode >= 200 && statusCode < 400 {
			m.service.HttpSuccessfulRequests.WithLabelValues().Inc()
			log.Println("Incremented successful requests counter")
		} else if statusCode >= 400 && statusCode < 600 {
			m.service.HttpUnsuccessfulRequests.WithLabelValues().Inc()
			log.Println("Incremented unsuccessful requests counter")
		}

		m.service.HttpTotalRequests.WithLabelValues().Inc()
		log.Println("Incremented total requests counter")
		m.service.AverageRequestDuration.WithLabelValues(r.Method, path).Set(duration)
		log.Printf("Set average request duration for %s %s", r.Method, path)
		m.service.RequestsPerTimeUnit.WithLabelValues(r.Method, path, "seconds").Inc()
		log.Printf("Incremented requests per time unit for %s %s", r.Method, path)
	})
}

func (m *Metrics) MetricsHandler() http.Handler {
	log.Println("Metrics handler invoked")
	return promhttp.HandlerFor(m.service.Registry, promhttp.HandlerOpts{})
}

func AdaptPrometheusHandler(handler http.Handler, metrics *Metrics) http.Handler {
	log.Println("Adapting handler with Prometheus metrics")
	return metrics.Count(handler)
}
