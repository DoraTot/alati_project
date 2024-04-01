package repositories

import (
	"errors"
	"projekat/model"
)

type ConfigInMemRepository struct {
	configs map[string]model.Config
}

func (c ConfigInMemRepository) GetConfig(name string, version float32) (*model.Config, error) {
	config, ok := c.configs[name]
	if !ok {
		return nil, errors.New("configuration not found")
	}

	if config.Version != version {
		return nil, errors.New("configuration version mismatch")
	}

	return &config, nil

}

func (c ConfigInMemRepository) AddConfig(config *model.Config) error {
	c.configs[config.Name] = *config
	return nil
}

func (c ConfigInMemRepository) DeleteConfig(name string) error {
	_, ok := c.configs[name]
	if !ok {
		return errors.New("configuration not found")
	}
	delete(c.configs, name)
	return nil
}

// todo: dodaj implementaciju metoda iz interfejsa ConfigRepository

func NewConfigInMemRepository() model.ConfigRepository {
	return ConfigInMemRepository{
		configs: make(map[string]model.Config),
	}
}
