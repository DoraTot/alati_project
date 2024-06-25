package services

import (
	"context"
	"fmt"
	"projekat/model"
)

type ConfigGroupService struct {
	repo model.ConfigGroupRepository
}

func NewConfigGroupService(repo model.ConfigGroupRepository) ConfigGroupService {
	return ConfigGroupService{
		repo: repo,
	}
}

func (s ConfigGroupService) Hello() {
	fmt.Println("hello from config group service")
}

func (s ConfigGroupService) AddConfigGroup(name string, version float32, configurations []model.ConfigForGroup, ctx context.Context) error {
	config := model.NewConfigGroup(name, version, configurations)
	return s.repo.AddConfigGroup(config, ctx)
}

func (s ConfigGroupService) GetConfigGroup(name string, version float32, ctx context.Context) (*model.ConfigGroup, error) {
	return s.repo.GetConfigGroup(name, version, ctx)
}

func (s ConfigGroupService) DeleteConfigGroup(name string, version float32, ctx context.Context) error {
	return s.repo.DeleteConfigGroup(name, version, ctx)
}
