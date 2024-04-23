package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"projekat/handlers"
	"projekat/model"
	"projekat/repositories"
	"projekat/services"
	"syscall"
	"time"
)

func main() {
	//repo := repositories.NewConfigConsulRepository() // Ovo koristiti kad budemo radili sa bazom
	repo2 := repositories.NewConfigGroupInMemRepository()
	repo := repositories.NewConfigInMemRepository(repo2)
	service := services.NewConfigService(repo)
	service.Hello()
	params := make(map[string]string)
	params["username"] = "pera"
	params["password"] = "pera"
	configs := model.NewConfig("db_config", 2.0, params)
	err := service.AddConfig(configs.Name, configs.Version, configs.Parameters)
	if err != nil {
		return
	}

	//For testing
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

	service2 := services.NewConfigGroupService(repo2)
	server2 := handlers.NewConfigGroupHandler(service2)

	router.HandleFunc("/config/", server.CreatePostHandler).Methods("POST")
	router.HandleFunc("/config/{name}/{version}/", server.Get).Methods("GET")
	router.HandleFunc("/config/{name}/{version}/", server.DelPostHandler).Methods("DELETE")
	router.HandleFunc("/configGroup/", server2.CreateConfigGroup).Methods("POST")
	router.HandleFunc("/configGroup/{name}/{version}/", server2.GetConfigGroup).Methods("GET")
	router.HandleFunc("/configGroup/{name}/{version}/", server2.DeleteConfigGroup).Methods("DELETE")
	router.HandleFunc("/config/{name}/{version}/{groupName}/{groupVersion}/", server.AddToConfigGroup).Methods("POST")
	router.HandleFunc("/config/{name}/{version}/{groupName}/{groupVersion}/", server.DeleteFromConfigGroup).Methods("DELETE")

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
