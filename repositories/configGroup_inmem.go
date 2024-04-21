package repositories

import (
	"errors"
	"fmt"
	"projekat/model"
)

type ConfigGroupInMemRepository struct {
	configs map[string]model.ConfigGroup
}

func (c ConfigGroupInMemRepository) GetConfigGroup(name string, version float32) (model.ConfigGroup, error) {
	key := fmt.Sprintf("%s/%d", name, version)
	config, ok := c.configs[key]
	if !ok {
		return model.ConfigGroup{}, errors.New("config not found")
	}
	return config, nil

}

func (c ConfigGroupInMemRepository) AddConfigGroup(config *model.ConfigGroup) error {
	key := fmt.Sprint("%s/%d", config.Name, config.Version)
	c.configs[key] = *config
	return nil
}

func (c ConfigGroupInMemRepository) DeleteConfigGroup(name string, version float32) error {
	key := fmt.Sprint("%s/%d", name, version)
	_, err := c.configs[key]
	if !err {
		return errors.New("configuration group does not exist")
	}
	delete(c.configs, key)
	return nil
}

// todo: dodaj implementaciju metoda iz interfejsa ConfigRepository

func NewConfigGroupInMemRepository() ConfigGroupInMemRepository {
	return ConfigGroupInMemRepository{
		configs: make(map[string]model.ConfigGroup),
	}
}
