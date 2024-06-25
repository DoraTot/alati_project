package services

import (
	"context"
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

func (s ConfigService) AddConfig(name string, version float32, parameters map[string]string, ctx context.Context) error {
	config := model.NewConfig(name, version, parameters)
	return s.repo.AddConfig(config, ctx)
}

func (s ConfigService) GetConfig(name string, version float32, ctx context.Context) (*model.Config, error) {
	return s.repo.GetConfig(name, version, ctx)
}

func (s ConfigService) DeleteConfig(name string, version float32, ctx context.Context) error {
	return s.repo.DeleteConfig(name, version, ctx)
}

// todo: implementiraj metode za dodavanje, brisanje, dobavljanje itd.
