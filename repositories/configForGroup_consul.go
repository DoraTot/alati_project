package repositories

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
	"log"
	"os"
	"projekat/model"
)

type ConfigForGroupConsulRepository struct {
	cli    *api.Client
	logger *log.Logger
}

func NewCFG(logger *log.Logger) (*ConfigForGroupConsulRepository, error) {
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
	return &ConfigForGroupConsulRepository{cli: client, logger: logger}, nil
}

func labelsMatch1(configLabels map[string]string, targetLabels map[string]string) bool {
	for key, value := range targetLabels {
		configValue, ok := configLabels[key]
		if !ok || configValue != value {
			return false
		}
	}
	return true
}

func (c ConfigForGroupConsulRepository) GetConfigsByLabels(groupName string, groupVersion float32, labels map[string]string) ([]model.ConfigForGroup, error) {
	if c.cli == nil {
		return nil, fmt.Errorf("Consul client is not initialized")
	}
	kv := c.cli.KV()
	groupKey := constructKeyForGroup(groupName, groupVersion)
	pair, _, err := kv.Get(groupKey, nil)
	if err != nil {
		return nil, err
	}
	if pair == nil {
		return nil, fmt.Errorf("configuration group '%s' with version %.2f does not exist", groupName, groupVersion)
	}
	var group model.ConfigGroup
	err = json.Unmarshal(pair.Value, &group)
	if err != nil {
		return nil, err
	}
	var matchingConfigs []model.ConfigForGroup
	for _, config := range group.Configurations {
		if labelsMatch1(config.Labels, labels) {
			matchingConfigs = append(matchingConfigs, config)
		}
	}
	return matchingConfigs, nil
}

func (c ConfigForGroupConsulRepository) DeleteConfigsByLabels(groupName string, groupVersion float32, labels map[string]string) error {
	kv := c.cli.KV()
	groupKey := constructKeyForGroup(groupName, groupVersion)
	pair, _, err := kv.Get(groupKey, nil)
	if err != nil {
		return err
	}
	if pair == nil {
		return fmt.Errorf("configuration group '%s' with version %.2f does not exist", groupName, groupVersion)
	}
	var group model.ConfigGroup
	err = json.Unmarshal(pair.Value, &group)
	if err != nil {
		return err
	}
	labelsFound := false
	for i := len(group.Configurations) - 1; i >= 0; i-- {
		config := group.Configurations[i]

		if labelsMatch1(config.Labels, labels) {
			group.Configurations = append(group.Configurations[:i], group.Configurations[i+1:]...)
			labelsFound = true
		}
	}
	if !labelsFound {
		return errors.New("labels not found")
	}
	updatedGroupJSON, err := json.Marshal(group)
	if err != nil {
		return err
	}
	p := &api.KVPair{Key: groupKey, Value: updatedGroupJSON}
	_, err = kv.Put(p, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c ConfigForGroupConsulRepository) AddToConfigGroup(config *model.ConfigForGroup, groupName string, groupVersion float32) error {

	kv := c.cli.KV()
	groupKey := constructKeyForGroup(groupName, groupVersion)

	pair, _, err := kv.Get(groupKey, nil)
	if err != nil {
		return err
	}
	if pair == nil {
		return fmt.Errorf("configuration group '%s' with version %.2f does not exist", groupName, groupVersion)
	}

	var group model.ConfigGroup
	err = json.Unmarshal(pair.Value, &group)
	if err != nil {
		return err
	}

	configForGroup := &model.ConfigForGroup{
		Name:       config.Name,
		Labels:     config.Labels,
		Parameters: config.Parameters,
	}

	group.Configurations = append(group.Configurations, *configForGroup)

	updatedGroupJSON, err := json.Marshal(group)
	if err != nil {
		return err
	}
	c.logger.Printf("Adding config to config group with SID: %s, Data: %s\n", groupKey, string(updatedGroupJSON))

	p := &api.KVPair{Key: groupKey, Value: updatedGroupJSON}
	_, err = kv.Put(p, nil)
	if err != nil {
		return err
	}
	c.logger.Println("Config successfully added to config group Consul KV:", groupKey)

	return nil
}

func (c ConfigForGroupConsulRepository) DeleteFromConfigGroup(configForGroupName string, groupName string, groupVersion float32) error {
	kv := c.cli.KV()

	groupKey := constructKeyForGroup(groupName, groupVersion)

	pair, _, err := kv.Get(groupKey, nil)
	if err != nil {
		return err
	}
	if pair == nil {
		return fmt.Errorf("configuration group '%s' with version %.2f does not exist", groupName, groupVersion)
	}

	var group model.ConfigGroup
	err = json.Unmarshal(pair.Value, &group)
	if err != nil {
		return err
	}

	found := false
	index := -1

	for i, configFromGroup := range group.Configurations {
		if configFromGroup.Name == configForGroupName {
			index = i
			found = true
			break
		}
	}

	if found {
		group.Configurations = append(group.Configurations[:index], group.Configurations[index+1:]...)
		c.logger.Println("Config successfully deleted from group Consul:")

		updatedGroupJSON, err := json.Marshal(group)
		if err != nil {
			return err
		}

		p := &api.KVPair{Key: groupKey, Value: updatedGroupJSON}
		_, err = kv.Put(p, nil)
		if err != nil {
			return err
		}

		return nil
	}

	return errors.New("configuration not found in the specified group")
}

func NewConfigForGroupConsulRepository() model.ConfigForGroupRepository {
	return ConfigForGroupConsulRepository{}
}
