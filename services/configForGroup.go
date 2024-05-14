package services

import "projekat/model"

type ConfigForGroupService struct {
	repo model.ConfigForGroupRepository
}

func NewConfigForGroupService(repo model.ConfigForGroupRepository) ConfigForGroupService {
	return ConfigForGroupService{
		repo: repo,
	}
}

func (s ConfigForGroupService) AddToConfigGroup(name string, labels map[string]string, parameters map[string]string, groupName string, groupVersion float32) error {
	config := model.NewConfigForGroup(name, labels, parameters)
	return s.repo.AddToConfigGroup(config, groupName, groupVersion)
}

func (s ConfigForGroupService) DeleteFromConfigGroup(configForGroupName string, groupName string, groupVersion float32) error {
	return s.repo.DeleteFromConfigGroup(configForGroupName, groupName, groupVersion)

}

func (s ConfigForGroupService) GetConfigsByLabels(groupName string, groupVersion float32, labels map[string]string) ([]model.ConfigForGroup, error) {
	return s.repo.GetConfigsByLabels(groupName, groupVersion, labels)
}
