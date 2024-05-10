package main

import (
	"projekat/repositories"
	"projekat/services"
)

func main() {
	repo := repositories.NewConfigConsulRepository()
	service := services.NewConfigService(repo)
	service.Hello()
	
}
