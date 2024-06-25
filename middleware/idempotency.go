package middleware

import (
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"projekat/model"
	"projekat/services"
	"sync"
)

type Idempotency struct {
	mux     sync.Mutex
	service *services.IdempotencyService
	Tracer  trace.Tracer
}

func NewIdempotency(idempotencyService *services.IdempotencyService, tracer trace.Tracer) *Idempotency {
	return &Idempotency{
		service: idempotencyService,
		Tracer:  tracer,
	}
}

func AdaptIdempotencyHandler(handler http.Handler, idempotencyMiddleware *Idempotency) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx, span := idempotencyMiddleware.Tracer.Start(r.Context(), "IdempotencyMiddleware.AdaptIdempotencyHandler")
		defer span.End()
		idempotencyMiddleware.mux.Lock()
		defer idempotencyMiddleware.mux.Unlock()

		if r.Method == http.MethodPost {
			idempotencyKey := r.Header.Get("Idempotency-Key")
			newRequest := model.IdempotencyRequest{}
			newRequest.SetKey(idempotencyKey)

			if idempotencyKey == "" {
				span.SetStatus(codes.Unset, "Key missing")
				http.Error(w, "Idempotency-Key header is missing", http.StatusBadRequest)
				return
			}

			processed, err := idempotencyMiddleware.service.Get(idempotencyKey, ctx)
			if err != nil {
				span.SetStatus(codes.Error, err.Error())
				http.Error(w, "Error checking idempotency: "+err.Error(), http.StatusInternalServerError)
				return
			}

			if processed {
				span.SetStatus(codes.Ok, "")
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte("Request already sent."))
				return
			}

			idempotencyMiddleware.service.Add(&newRequest, ctx)
			handler.ServeHTTP(w, r)
			return
		}

		handler.ServeHTTP(w, r)
		span.SetStatus(codes.Ok, "")
	})
}
