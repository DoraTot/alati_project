// Post API
//
//	Title: Post API
//
//	Schemes: http
//	Version: 0.0.1
//	BasePath: /
//
//	Produces:
//	  - application/json
//
// swagger:meta

package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"os"
	"os/signal"
	"projekat/configuration"
	"projekat/handlers"
	middleware2 "projekat/middleware"
	"projekat/model"
	"projekat/repositories"
	"projekat/services"
	"syscall"
	"time"
)

func main() {

	dbHost := os.Getenv("DB")
	dbPort := os.Getenv("DBPORT")

	if dbHost == "" || dbPort == "" {
		log.Fatal("DB and DBPORT environment variables must be set")
	}

	// Set Consul address
	consulAddress := "http://" + dbHost + ":" + dbPort
	os.Setenv("CONSUL_HTTP_ADDR", consulAddress)

	port := os.Getenv("PORT") // set port for consul
	if len(port) == 0 {
		port = "8000"
	}
	logger := log.New(os.Stdout, "[config-api] ", log.LstdFlags)

	// For Tracing
	cfg := configuration.GetConfiguration()
	ctx := context.Background()
	exp, err := newExporter(cfg.JaegerAddress)
	if err != nil {
		log.Fatalf("failed to initialize exporter: %v", err)
	}
	tp := newTraceProvider(exp)
	defer func() { _ = tp.Shutdown(ctx) }()
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	tracer := tp.Tracer("alati_projekat")

	repo, err := repositories.New(logger, tracer) // new consul repo for configs
	if err != nil {
		logger.Fatal("Failed to create repository:", err)
	}

	repoCG, err2 := repositories.NewCG(logger, tracer) // consul for configGroup
	if err2 != nil {
		logger.Fatal("Failed to create repository for configGroup:", err)
	}

	repoCFG, err3 := repositories.NewCFG(logger, tracer)
	if err3 != nil {
		logger.Fatal("Failed to create repository for configForGroup:", err3)
	}

	//repo2 := repositories.NewConfigGroupInMemRepository()

	service := services.NewConfigService(repo)
	service1 := services.NewConfigForGroupService(repoCFG)
	service.Hello()
	params := make(map[string]string)
	params["username"] = "pera"
	params["password"] = "pera"
	configs := model.NewConfig("db_config", 2.0, params)
	err = service.AddConfig(configs.Name, configs.Version, configs.Parameters, ctx)
	if err != nil {
		return
	}

	var limiter = rate.NewLimiter(0.167, 10) //For testing
	name := "db_config"
	version := float32(2.0)
	config, err := repo.GetConfig(name, version, ctx)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Config:", config)
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	idempotencyService := services.NewIdempotencyService(*repo, tracer)
	idempotencyMiddleware := middleware2.NewIdempotency(&idempotencyService, tracer)
	metricsService := services.NewMetricsService()
	metricsMiddleware := middleware2.NewMetrics(metricsService)

	router := mux.NewRouter()
	router.StrictSlash(true)
	router.Use(otelmux.Middleware("alati_projekat"))

	router.Use(func(next http.Handler) http.Handler {
		return middleware2.AdaptIdempotencyHandler(next, idempotencyMiddleware)
	})

	router.Use(func(next http.Handler) http.Handler {
		return middleware2.AdaptPrometheusHandler(next, metricsMiddleware)
	})

	server := handlers.NewConfigHandler(logger, service, tracer)

	server1 := handlers.NewConfigForGroupHandler(service1, tracer)
	service2 := services.NewConfigGroupService(repoCG)
	server2 := handlers.NewConfigGroupHandler(service2, tracer)

	router.Handle("/config/", middleware2.RateLimit(limiter, server.CreatePostHandler)).Methods("POST")
	router.Handle("/config/{name}/{version}/", middleware2.RateLimit(limiter, server.Get)).Methods("GET")
	router.Handle("/config/{name}/{version}/", middleware2.RateLimit(limiter, server.DelPostHandler)).Methods("DELETE")
	router.Handle("/configGroup/", middleware2.RateLimit(limiter, server2.CreateConfigGroup)).Methods("POST")
	router.Handle("/configGroup/{name}/{version}/", middleware2.RateLimit(limiter, server2.GetConfigGroup)).Methods("GET")
	router.Handle("/configGroup/{name}/{version}/", middleware2.RateLimit(limiter, server2.DeleteConfigGroup)).Methods("DELETE")
	router.Handle("/config/configGroup/{groupName}/{groupVersion}/", middleware2.RateLimit(limiter, server1.AddToConfigGroup)).Methods("POST")
	router.Handle("/config/{name}/{groupName}/{groupVersion}/", middleware2.RateLimit(limiter, server1.DeleteFromConfigGroup)).Methods("DELETE")

	router.Handle("/configGroup/{groupName}/{groupVersion}/{labels}", middleware2.RateLimit(limiter, server1.DeleteConfigsByLabels)).Methods("DELETE")
	router.Handle("/configGroup/{groupName}/{groupVersion}/{labels}", middleware2.RateLimit(limiter, server1.GetConfigsByLabels)).Methods("GET")
	//router.HandleFunc("/swagger.yaml", middleware2.SwaggerHandler).Methods("GET")
	//router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./"))))

	router.HandleFunc("/swagger.yaml", middleware2.SwaggerHandler).Methods("GET")

	// SwaggerUI
	optionsDevelopers := middleware.SwaggerUIOpts{SpecURL: "http://localhost:8081/swagger.yaml"} // Note the use of http://localhost:8081
	developerDocumentationHandler := middleware.SwaggerUI(optionsDevelopers, nil)
	router.Handle("/docs", developerDocumentationHandler)

	router.Handle("/metrics", metricsMiddleware.MetricsHandler()).Methods("GET")

	// CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "PUT", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	// start server
	srv := &http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: corsHandler.Handler(router),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Fatal(err)
			}
		}
	}()

	<-quit
	log.Println("Service shutting down...")

	// gracefully stop server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("Server stopped")
}

func newExporter(address string) (*jaeger.Exporter, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(address)))
	if err != nil {
		return nil, err
	}
	return exp, nil
}

func newTraceProvider(exp sdktrace.SpanExporter) *sdktrace.TracerProvider {
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("alati_projekat"),
		),
	)
	if err != nil {
		log.Fatalf("failed to create resource: %v", err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)
}
