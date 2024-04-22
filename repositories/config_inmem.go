package repositories

import (
	"errors"
	"fmt"
	"projekat/model"
)

type ConfigInMemRepository struct {
	configs      map[string]model.Config
	configGroups map[string][]string
}

func (c ConfigInMemRepository) GetConfig(name string, version float32) (*model.Config, error) {
	key := fmt.Sprintf("%s/%d", name, version)
	config, ok := c.configs[key]
	if !ok {
		return &model.Config{}, errors.New("config not found")
	}
	return &config, nil

}

func (c ConfigInMemRepository) AddConfig(config *model.Config) error {

	key := fmt.Sprint("%s/%d", config.Name, config.Version)
	c.configs[key] = *config
	return nil
}

//	func (c ConfigInMemRepository) DeleteConfig(name string, version float32) error {
//		key := fmt.Sprint("%s/%d", name, version)
//		_, err := c.configs[key]
//		if !err {
//			return errors.New("configuration does not exist")
//		}
//		delete(c.configs, key)
//		return nil
//	}
func (c ConfigInMemRepository) DeleteConfig(name string) error {
	key := fmt.Sprint("%s/%d", name)
	_, err := c.configs[key]
	if !err {
		return errors.New("configuration does not exist")
	}
	delete(c.configs, key)
	return nil
}

func (c ConfigInMemRepository) AddToConfigGroup(config *model.Config, groupName string) error {
	key := fmt.Sprintf("%s/%f", config.Name, config.Version)
	_, ok := c.configs[key]
	if !ok {
		return errors.New("configuration does not exist")
	}

	// Add the configuration to the specified group
	if _, groupExists := c.configGroups[groupName]; !groupExists {
		c.configGroups[groupName] = []string{key}
	} else {
		c.configGroups[groupName] = append(c.configGroups[groupName], key)
	}
	return nil
}

func (c ConfigInMemRepository) DeleteFromConfigGroup(config *model.Config, groupName string) error {
	groupConfigs, ok := c.configGroups[groupName]
	if !ok {
		return errors.New("group does not exist")
	}

	// Find and remove the configuration from the group
	key := fmt.Sprintf("%s/%f", config.Name, config.Version)
	found := false
	for i, configKey := range groupConfigs {
		if configKey == key {
			// Remove the config key from the group slice
			c.configGroups[groupName] = append(groupConfigs[:i], groupConfigs[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		return errors.New("configuration not found in the specified group")
	}
	return nil
}

// todo: dodaj implementaciju metoda iz interfejsa ConfigRepository

func NewConfigInMemRepository() ConfigInMemRepository {
	return ConfigInMemRepository{
		configs:      make(map[string]model.Config),
		configGroups: make(map[string][]string),
	}
}
