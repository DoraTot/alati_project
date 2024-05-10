package repositories

import (
	"errors"
	"fmt"
	"projekat/model"
)

type ConfigInMemRepository struct {
	Configs      map[string]model.Config
	ConfigGroups *ConfigGroupInMemRepository
}

func (c ConfigInMemRepository) GetConfig(name string, version float32) (*model.Config, error) {
	key := fmt.Sprintf("%s/%.2f", name, version)
	config, ok := c.Configs[key]
	if !ok {
		return &model.Config{}, errors.New("config not found")
	}
	return &config, nil

}

func (c ConfigInMemRepository) AddConfig(config *model.Config) error {

	key := fmt.Sprintf("%s/%.2f", config.Name, config.Version)
	c.Configs[key] = *config
	return nil
}

func (c ConfigInMemRepository) DeleteConfig(name string, version float32) error {
	key := fmt.Sprintf("%s/%.2f", name, version)
	if _, ok := c.Configs[key]; !ok {
		return errors.New("configuration does not exist")
	}
	delete(c.Configs, key)
	fmt.Printf("Deleting configuration: %s\n", key)
	return nil
}

func (c ConfigInMemRepository) AddToConfigGroup(config *model.Config, groupName string, groupVersion float32) error {
	key := fmt.Sprintf("%s/%.2f", config.Name, config.Version)
	_, ok := c.Configs[key]
	if !ok {
		return errors.New("configuration does not exist")
	}

	group, err := c.ConfigGroups.GetConfigGroup(groupName, groupVersion)
	if err != nil {
		return err // Error fetching the group
	}
	if group == nil {
		return fmt.Errorf("configuration group '%s' with version %.2f does not exist", groupName, groupVersion)
	}
	group.Configurations = append(group.Configurations, *config)

	return nil
}

func (c ConfigInMemRepository) DeleteFromConfigGroup(config *model.Config, groupName string, groupVersion float32) error {
	key := fmt.Sprintf("%s/%.2f", config.Name, config.Version)
	_, ok := c.Configs[key]
	if !ok {
		return errors.New("configuration does not exist")
	}

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
		if configFromGroup.Name == config.Name && configFromGroup.Version == config.Version {
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

// todo: dodaj implementaciju metoda iz interfejsa ConfigRepository

func NewConfigInMemRepository(groupRepo *ConfigGroupInMemRepository) *ConfigInMemRepository {
	return &ConfigInMemRepository{
		Configs:      make(map[string]model.Config),
		ConfigGroups: groupRepo, // Assuming ConfigGroups should point to an instance of ConfigGroupInMemRepository
	}
}
