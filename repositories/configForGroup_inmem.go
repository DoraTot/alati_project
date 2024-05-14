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

func (c ConfigForGroupInMemRepository) GetConfigsByLabels(groupName string, groupVersion float32, labels map[string]string) ([]model.ConfigForGroup, error) {
	group, err := c.ConfigGroups.GetConfigGroup(groupName, groupVersion)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, fmt.Errorf("configuration group '%s' with version %.2f does not exist", groupName, groupVersion)
	}

	var matchingConfigs []model.ConfigForGroup

	// Iterate through configurations in the group
	for _, config := range group.Configurations {
		// Check if the labels match
		if labelsMatch(config.Labels, labels) {
			// If labels match, add the configuration to the result
			matchingConfigs = append(matchingConfigs, config)
		}
	}
	fmt.Println(matchingConfigs)
	return matchingConfigs, nil
}

func labelsMatch(configLabels map[string]string, targetLabels map[string]string) bool {
	// Iterate through targetLabels and check if each key-value pair exists in configLabels
	for key, value := range targetLabels {
		// If the key doesn't exist in the configuration labels or the values don't match, return false
		configValue, ok := configLabels[key]
		if !ok || configValue != value {
			return false
		}
	}
	// If all target labels are present and match, return true
	return true
}

func (c ConfigForGroupInMemRepository) DeleteConfigsByLabels(groupName string, groupVersion float32, labels map[string]string) error {
	group, err := c.ConfigGroups.GetConfigGroup(groupName, groupVersion)
	if err != nil {
		return err
	}
	if group == nil {
		return fmt.Errorf("configuration group '%s' with version %.2f does not exist", groupName, groupVersion)
	}

	labelsFound := false
	indicesToDelete := []int{}
	for i, config := range group.Configurations {
		if labelsMatch(config.Labels, labels) {
			indicesToDelete = append(indicesToDelete, i)
			labelsFound = true
		}
	}

	// Delete configurations starting from the end to avoid index shifting
	for i := len(indicesToDelete) - 1; i >= 0; i-- {
		index := indicesToDelete[i]
		group.Configurations = append(group.Configurations[:index], group.Configurations[index+1:]...)
	}
	if labelsFound == true {
		return nil
	}
	return errors.New("labels not found")
}

func NewConfigForGroupInMemRepository(groupRepo *ConfigGroupInMemRepository) *ConfigForGroupInMemRepository {
	return &ConfigForGroupInMemRepository{
		Configs:      make(map[string]model.ConfigForGroup),
		ConfigGroups: groupRepo, // Assuming ConfigGroups should point to an instance of ConfigGroupInMemRepository
	}
}
