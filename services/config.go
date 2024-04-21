package services

import (
	"fmt"
	"projekat/model"
)

type ConfigService struct {
	repo model.ConfigRepository
}

func NewConfigService(repo model.ConfigRepository) ConfigService {
	return ConfigService{
		repo: repo,
	}
}

func (s ConfigService) Hello() {
	fmt.Println("hello from config service")
}

func (s ConfigService) AddConfig(name string, version float32, parameters map[string]string) error {
	config := model.NewConfig(name, version, parameters)
	return s.repo.AddConfig(config)
}

func (s ConfigService) GetConfig(name string, version float32) (*model.Config, error) {
	return s.repo.GetConfig(name, version)
}

func (s ConfigService) DeleteConfig(name string, version float32) error {
	return s.repo.DeleteConfig(name)
}

// todo: implementiraj metode za dodavanje, brisanje, dobavljanje itd.
