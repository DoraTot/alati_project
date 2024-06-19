package middleware

import (
	"net/http"
	"projekat/model"
	"projekat/services"
	"sync"
)

type Idempotency struct {
	mux     sync.Mutex
	service *services.IdempotencyService
}

func NewIdempotency(idempotencyService *services.IdempotencyService) *Idempotency {
	return &Idempotency{
		service: idempotencyService,
	}
}

func AdaptIdempotencyHandler(handler http.Handler, idempotencyMiddleware *Idempotency) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		idempotencyMiddleware.mux.Lock()
		defer idempotencyMiddleware.mux.Unlock()

		if r.Method == http.MethodPost {
			idempotencyKey := r.Header.Get("Idempotency-Key")
			newRequest := model.IdempotencyRequest{}
			newRequest.SetKey(idempotencyKey)

			if idempotencyKey == "" {
				http.Error(w, "Idempotency-Key header is missing", http.StatusBadRequest)
				return
			}

			processed, err := idempotencyMiddleware.service.Get(idempotencyKey)
			if err != nil {
				http.Error(w, "Error checking idempotency: "+err.Error(), http.StatusInternalServerError)
				return
			}

			if processed {
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte("Request already sent."))
				return
			}

			idempotencyMiddleware.service.Add(&newRequest)
			handler.ServeHTTP(w, r)
			return
		}

		handler.ServeHTTP(w, r)
	})
}
