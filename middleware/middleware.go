package middleware

import (
	"encoding/json"
	"golang.org/x/time/rate"
	"net/http"
)

func RateLimit(limiter *rate.Limiter, next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			message := map[string]string{
				"message": "Rate limit exceeded, try again later!",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			err := json.NewEncoder(w).Encode(message)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			next(w, r)
		}
	})
}

func SwaggerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		return
	}

	http.ServeFile(w, r, "./swagger.yaml")
}
