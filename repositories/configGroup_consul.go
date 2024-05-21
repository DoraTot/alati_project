package repositories

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul/api"
	"log"
	"os"
	"projekat/model"
)

type ConfigGroupConsulRepository struct {
	cli    *api.Client
	logger *log.Logger
}

func NewCG(logger *log.Logger) (*ConfigGroupConsulRepository, error) {
	db := os.Getenv("DB")
	dbport := os.Getenv("DBPORT")
	if db == "" || dbport == "" {
		return nil, fmt.Errorf("environment variables DB and DBPORT must be set")
	}
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%s", db, dbport)
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &ConfigGroupConsulRepository{cli: client, logger: logger}, nil
}

func (c ConfigGroupConsulRepository) GetConfigGroup(name string, version float32) (*model.ConfigGroup, error) {
	kv := c.cli.KV()
	key := constructKeyForGroup(name, version)
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return nil, err
	}
	if pair == nil {
		return nil, fmt.Errorf("configuration group '%s' with version %.1f not found", name, version)
	}

	configGroup := &model.ConfigGroup{}
	err = json.Unmarshal(pair.Value, configGroup)
	if err != nil {
		return nil, err
	}
	return configGroup, nil

}

func (c ConfigGroupConsulRepository) AddConfigGroup(config *model.ConfigGroup) error {
	kv := c.cli.KV()
	key := constructKeyForGroup(config.Name, config.Version)

	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	c.logger.Printf("Adding config group with SID: %s, Data: %s\n", key, string(data))

	p := &api.KVPair{Key: key, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		c.logger.Println("Error adding config group:", err)
		return err
	}
	c.logger.Println("Config group added successfully", key)
	return nil

}

func (c ConfigGroupConsulRepository) DeleteConfigGroup(name string, version float32) error {
	kv := c.cli.KV()
	_, err := kv.Delete(constructKeyForGroup(name, version), nil)
	if err != nil {
		c.logger.Println("Error deleting config group:", err)
		return err
	}

	c.logger.Println("Config group deleted successfully", constructKeyForGroup(name, version))
	return nil
}

func NewConfigGroupConsulRepository() model.ConfigGroupRepository {
	return ConfigGroupConsulRepository{}
}
