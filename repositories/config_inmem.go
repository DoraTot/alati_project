package repositories

import (
	"errors"
	"fmt"
	"projekat/model"
)

type ConfigInMemRepository struct {
	Configs map[string]model.Config
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

// todo: dodaj implementaciju metoda iz interfejsa ConfigRepository

func NewConfigInMemRepository() *ConfigInMemRepository {
	return &ConfigInMemRepository{
		Configs: make(map[string]model.Config),
	}
}
