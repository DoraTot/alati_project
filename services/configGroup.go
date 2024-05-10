package services

import (
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

func (s ConfigGroupService) AddConfigGroup(name string, version float32, configurations []model.Config) error {
	config := model.NewConfigGroup(name, version, configurations)
	return s.repo.AddConfigGroup(config)
}

func (s ConfigGroupService) GetConfigGroup(name string, version float32) (*model.ConfigGroup, error) {
	return s.repo.GetConfigGroup(name, version)
}

func (s ConfigGroupService) DeleteConfigGroup(name string, version float32) error {
	return s.repo.DeleteConfigGroup(name, version)
}
