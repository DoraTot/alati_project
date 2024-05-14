package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"os"
	"os/signal"
	"projekat/handlers"
	"projekat/middleware"
	"projekat/model"
	"projekat/repositories"
	"projekat/services"
	"syscall"
	"time"
)

func main() {
	//repo := repositories.NewConfigConsulRepository() // Ovo koristiti kad budemo radili sa bazom

	repo2 := repositories.NewConfigGroupInMemRepository()
	repo := repositories.NewConfigInMemRepository()
	repo1 := repositories.NewConfigForGroupInMemRepository(repo2)

	service := services.NewConfigService(repo)
	service1 := services.NewConfigForGroupService(repo1)
	service.Hello()
	params := make(map[string]string)
	params["username"] = "pera"
	params["password"] = "pera"
	configs := model.NewConfig("db_config", 2.0, params)
	err := service.AddConfig(configs.Name, configs.Version, configs.Parameters)
	if err != nil {
		return
	}

	var limiter = rate.NewLimiter(0.167, 10) //For testing
	name := "db_config"
	version := float32(2.0)
	config, err := repo.GetConfig(name, version)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Config:", config)
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	router := mux.NewRouter()
	router.StrictSlash(true)

	server := handlers.NewConfigHandler(service)

	server1 := handlers.NewConfigForGroupHandler(service1)
	service2 := services.NewConfigGroupService(repo2)
	server2 := handlers.NewConfigGroupHandler(service2)

	router.Handle("/config/", middleware.RateLimit(limiter, server.CreatePostHandler)).Methods("POST")
	router.Handle("/config/{name}/{version}/", middleware.RateLimit(limiter, server.Get)).Methods("GET")
	router.Handle("/config/{name}/{version}/", middleware.RateLimit(limiter, server.DelPostHandler)).Methods("DELETE")
	router.Handle("/configGroup/", middleware.RateLimit(limiter, server2.CreateConfigGroup)).Methods("POST")
	router.Handle("/configGroup/{name}/{version}/", middleware.RateLimit(limiter, server2.GetConfigGroup)).Methods("GET")
	router.Handle("/configGroup/{name}/{version}/", middleware.RateLimit(limiter, server2.DeleteConfigGroup)).Methods("DELETE")
	router.Handle("/config/configGroup/", middleware.RateLimit(limiter, server1.AddToConfigGroup)).Methods("POST")
	router.Handle("/config/{name}/{groupName}/{groupVersion}/", middleware.RateLimit(limiter, server1.DeleteFromConfigGroup)).Methods("DELETE")

	//router.HandleFunc("/configGroup/{name}/{version}/{labels}", server1.DeleteConfigsByLabels).Methods("DELETE")
	router.HandleFunc("/configGroup/{groupName}/{groupVersion}/{labels}", server1.GetConfigsByLabels).Methods("GET")

	srv := &http.Server{
		Addr:    "0.0.0.0:8000",
		Handler: router,
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("Server stopped")
}
