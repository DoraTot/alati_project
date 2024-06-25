package services

import (
	"context"
	"projekat/model"
)

type ConfigForGroupService struct {
	repo model.ConfigForGroupRepository
}

func NewConfigForGroupService(repo model.ConfigForGroupRepository) ConfigForGroupService {
	return ConfigForGroupService{
		repo: repo,
	}
}

func (s ConfigForGroupService) AddToConfigGroup(name string, labels map[string]string, parameters map[string]string, groupName string, groupVersion float32, ctx context.Context) error {
	config := model.NewConfigForGroup(name, labels, parameters)
	return s.repo.AddToConfigGroup(config, groupName, groupVersion, ctx)
}

func (s ConfigForGroupService) DeleteFromConfigGroup(configForGroupName string, groupName string, groupVersion float32, ctx context.Context) error {
	return s.repo.DeleteFromConfigGroup(configForGroupName, groupName, groupVersion, ctx)

}

func (s ConfigForGroupService) GetConfigsByLabels(groupName string, groupVersion float32, labels map[string]string, ctx context.Context) ([]model.ConfigForGroup, error) {
	return s.repo.GetConfigsByLabels(groupName, groupVersion, labels, ctx)
}

func (s ConfigForGroupService) DeleteConfigsByLabels(groupName string, groupVersion float32, labels map[string]string, ctx context.Context) error {
	return s.repo.DeleteConfigsByLabels(groupName, groupVersion, labels, ctx)
}
