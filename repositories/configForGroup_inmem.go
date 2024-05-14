package repositories

import (
	"errors"
	"fmt"
	"projekat/model"
)

type ConfigForGroupInMemRepository struct {
	Configs      map[string]model.ConfigForGroup
	ConfigGroups *ConfigGroupInMemRepository
}

func (c ConfigForGroupInMemRepository) GetConfigForGroup(name string, labels string) (*model.ConfigForGroup, error) {
	key := fmt.Sprintf("%s/%s", name, labels)
	ConfigForGroup, ok := c.Configs[key]
	if !ok {
		return &model.ConfigForGroup{}, errors.New("config not found")
	}
	return &ConfigForGroup, nil

}

func (c ConfigForGroupInMemRepository) AddToConfigGroup(config *model.ConfigForGroup, groupName string, groupVersion float32) error {

	group, err := c.ConfigGroups.GetConfigGroup(groupName, groupVersion)
	if err != nil {
		return err // Error fetching the group
	}
	if group == nil {
		return fmt.Errorf("configuration group '%s' with version %.2f does not exist", groupName, groupVersion)
	}

	configForGroup := &model.ConfigForGroup{
		Name:       config.Name,
		Labels:     config.Labels,
		Parameters: config.Parameters,
	}

	group.Configurations = append(group.Configurations, *configForGroup)

	return nil
}

func (c ConfigForGroupInMemRepository) DeleteFromConfigGroup(configForGroupName string, groupName string, groupVersion float32) error {

	group, err := c.ConfigGroups.GetConfigGroup(groupName, groupVersion)
	if err != nil {
		return err // Error fetching the group
	}
	if group == nil {
		return fmt.Errorf("configuration group '%s' with version %.2f does not exist", groupName, groupVersion)
	}

	found := false
	index := -1
	for i, configFromGroup := range group.Configurations {
		// Check if the names match
		if configFromGroup.Name == configForGroupName {
			index = i
			found = true
			break

		}
	}

	if found {
		group.Configurations = append(group.Configurations[:index], group.Configurations[index+1:]...)
		return nil
	}

	return errors.New("configuration not found in the specified group")
}

func NewConfigForGroupInMemRepository(groupRepo *ConfigGroupInMemRepository) *ConfigForGroupInMemRepository {
	return &ConfigForGroupInMemRepository{
		Configs:      make(map[string]model.ConfigForGroup),
		ConfigGroups: groupRepo, // Assuming ConfigGroups should point to an instance of ConfigGroupInMemRepository
	}
}
