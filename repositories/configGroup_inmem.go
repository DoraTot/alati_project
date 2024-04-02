package repositories

import (
	"errors"
	"projekat/model"
)

type ConfigGroupInMemRepository struct {
	configs map[string]model.ConfigGroup
}

func (c ConfigGroupInMemRepository) GetConfigGroup(name string, version float32) (*model.ConfigGroup, error) {
	config, ok := c.configs[name]
	if !ok {
		return nil, errors.New("configuration group not found")
	}

	if config.Version != version {
		return nil, errors.New("configuration version mismatch")
	}

	return &config, nil

}

func (c ConfigGroupInMemRepository) AddConfigGroup(config *model.ConfigGroup) error {
	c.configs[config.Name] = *config
	return nil
}

func (c ConfigGroupInMemRepository) DeleteConfigGroup(name string, version float32) error {
	config, ok := c.configs[name]
	if !ok {
		return errors.New("configuration not found")
	}
	if config.Version != version {
		return errors.New("configuration version mismatch")
	}
	delete(c.configs, name)
	return nil
}

// todo: dodaj implementaciju metoda iz interfejsa ConfigRepository

func NewConfigGroupInMemRepository() model.ConfigGroupRepository {
	return ConfigGroupInMemRepository{
		configs: make(map[string]model.ConfigGroup),
	}
}
