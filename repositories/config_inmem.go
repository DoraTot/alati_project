package repositories

import (
	"errors"
	"fmt"
	"projekat/model"
)

type ConfigInMemRepository struct {
	configs map[string]model.Config
}

func (c ConfigInMemRepository) GetConfig(name string, version float32) (*model.Config, error) {
	key := fmt.Sprint("%s/%d", name, version)
	config, err := c.configs[key]
	if !err {
		return nil, errors.New("configuration not found")
	}

	if config.Version != version {
		return nil, errors.New("configuration version mismatch")
	}

	return &config, nil

}

func (c ConfigInMemRepository) AddConfig(config *model.Config) error {

	key := fmt.Sprint("%s/%d", config.Name, config.Version)
	c.configs[key] = *config
	return nil
}

func (c ConfigInMemRepository) DeleteConfig(name string, version float32) error {
	key := fmt.Sprint("%s/%d", name, version)
	_, err := c.configs[key]
	if !err {
		return errors.New("configuration does not exist")
	}
	delete(c.configs, key)
	return nil
}

// todo: dodaj implementaciju metoda iz interfejsa ConfigRepository

//func NewConfigInMemRepository() model.ConfigRepository {
//	return ConfigInMemRepository{
//		configs: make(map[string]model.Config),
//	}
//}
