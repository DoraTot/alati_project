package repositories

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul/api"
	"log"
	"os"
	"projekat/model"
)

type ConfigConsulRepository struct {
	cli    *api.Client
	logger *log.Logger
}

func New(logger *log.Logger) (*ConfigConsulRepository, error) {
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
	return &ConfigConsulRepository{cli: client, logger: logger}, nil
}

// swagger:route GET /config/{name}/{version}/ getConfig
// Get config
//
// responses:
//
//	200: ResponseConfig
func (c ConfigConsulRepository) GetConfig(name string, version float32) (*model.Config, error) {
	kv := c.cli.KV()
	key := constructKey(name, version)
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return nil, err
	}
	if pair == nil {
		return nil, fmt.Errorf("configuration '%s' with version %.1f not found", name, version)
	}
	config := &model.Config{}
	err = json.Unmarshal(pair.Value, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// swagger:route POST /config/ config addConfig
// Add new config
//
// responses:
//
//	415: ErrorResponse
//	400: ErrorResponse
//	201: ResponseConfig
func (c ConfigConsulRepository) AddConfig(config *model.Config) error {
	kv := c.cli.KV()
	key := constructKey(config.Name, config.Version)

	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	c.logger.Printf("Adding config with SID: %s, Data: %s\n", key, string(data))

	p := &api.KVPair{Key: key, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		c.logger.Println("Error putting config to Consul KV:", err)
		return err
	}
	c.logger.Println("Config successfully added to Consul KV:", key)
	return nil
}

// swagger:route DELETE /config/{name}/{version}/ config deleteConfig
// Delete config
//
// responses:
//
//	404: ErrorResponse
//	204: NoContentResponse
func (c ConfigConsulRepository) DeleteConfig(name string, version float32) error {
	kv := c.cli.KV()
	_, err := kv.Delete(constructKey(name, version), nil)
	if err != nil {
		return err
	}
	c.logger.Println("Config successfully deleted from Consul:", name)

	return nil
}

// todo: dodaj implementaciju metoda iz interfejsa ConfigRepository

func NewConfigConsulRepository() model.ConfigRepository {
	return ConfigConsulRepository{}
}
