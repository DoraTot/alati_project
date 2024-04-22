package repositories

import (
	"errors"
	"fmt"
	"projekat/model"
)

type ConfigGroupInMemRepository struct {
	Configs map[string]model.ConfigGroup
}

func (c ConfigGroupInMemRepository) GetConfigGroup(name string, version float32) (*model.ConfigGroup, error) {
	key := fmt.Sprintf("%s/%.2f", name, version)
	config, ok := c.Configs[key]
	if !ok {
		return nil, fmt.Errorf("configGroup '%s' with version %.2f not found", name, version)
	}
	return &config, nil

}

func (c ConfigGroupInMemRepository) AddConfigGroup(config *model.ConfigGroup) error {
	key := fmt.Sprintf("%s/%.2f", config.Name, config.Version)
	c.Configs[key] = *config
	return nil
}

func (c ConfigGroupInMemRepository) DeleteConfigGroup(name string, version float32) error {
	key := fmt.Sprint("%s/%d", name, version)
	_, err := c.Configs[key]
	if !err {
		return errors.New("configuration group does not exist")
	}
	delete(c.Configs, key)
	return nil
}

// todo: dodaj implementaciju metoda iz interfejsa ConfigRepository

func NewConfigGroupInMemRepository() ConfigGroupInMemRepository {
	return ConfigGroupInMemRepository{
		Configs: make(map[string]model.ConfigGroup),
	}
}
