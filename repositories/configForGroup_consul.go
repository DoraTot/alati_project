package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"log"
	"os"
	"projekat/model"
)

type ConfigForGroupConsulRepository struct {
	cli    *api.Client
	logger *log.Logger
	Tracer trace.Tracer
}

func NewCFG(logger *log.Logger, tracer trace.Tracer) (*ConfigForGroupConsulRepository, error) {
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
	return &ConfigForGroupConsulRepository{cli: client, logger: logger, Tracer: tracer}, nil
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

// swagger:route GET /configGroup/{groupName}/{groupVersion}/{labels} getConfigsByLabels
// Get all configForGroups by labels
//
// responses:
//
//	200: []ResponseConfigForGroup
func (c ConfigForGroupConsulRepository) GetConfigsByLabels(groupName string, groupVersion float32, labels map[string]string, ctx context.Context) ([]model.ConfigForGroup, error) {
	_, span := c.Tracer.Start(ctx, "ConfigForGroupConsulRepository.GetConfigsByLabels")
	defer span.End()

	if c.cli == nil {
		span.SetStatus(codes.Error, "Consul not working")
		return nil, fmt.Errorf("Consul client is not initialized")
	}
	kv := c.cli.KV()
	groupKey := constructKeyForGroup(groupName, groupVersion)
	pair, _, err := kv.Get(groupKey, nil)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	if pair == nil {
		span.SetStatus(codes.Error, "Pair not found")
		return nil, fmt.Errorf("configuration group '%s' with version %.2f does not exist", groupName, groupVersion)
	}
	var group model.ConfigGroup
	err = json.Unmarshal(pair.Value, &group)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	var matchingConfigs []model.ConfigForGroup
	for _, config := range group.Configurations {
		if labelsMatch1(config.Labels, labels) {
			matchingConfigs = append(matchingConfigs, config)
		}
	}
	span.SetStatus(codes.Ok, "Success getting configuration group")
	return matchingConfigs, nil
}

// swagger:route DELETE /configGroup/{groupName}/{groupVersion}/{labels} deleteConfigsByLabels
// Delete configForGroup by labels
//
// responses:
//
//	404: ErrorResponse
//	204: NoContentResponse
func (c ConfigForGroupConsulRepository) DeleteConfigsByLabels(groupName string, groupVersion float32, labels map[string]string, ctx context.Context) error {
	_, span := c.Tracer.Start(ctx, "ConfigForGroupConsulRepository.DeleteConfigsByLabels")
	defer span.End()

	kv := c.cli.KV()
	groupKey := constructKeyForGroup(groupName, groupVersion)
	pair, _, err := kv.Get(groupKey, nil)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
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
		span.SetStatus(codes.Error, "labels not found")
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
	span.SetStatus(codes.Ok, "Success deleting configuration group by labels")
	return nil
}

// swagger:route POST /config/configGroup/ addToConfigGroup
// Add config to group
//
// responses:
//
//	415: ErrorResponse
//	400: ErrorResponse
//	201: ResponseConfigForGroup
func (c ConfigForGroupConsulRepository) AddToConfigGroup(config *model.ConfigForGroup, groupName string, groupVersion float32, ctx context.Context) error {
	_, span := c.Tracer.Start(ctx, "ConfigForGroupConsulRepository.AddToConfigGroup")
	defer span.End()

	if c.cli == nil {
		err := errors.New("Consul client is nil")
		log.Printf("Error: %v", err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	kv := c.cli.KV()
	if kv == nil {
		err := errors.New("KV store is nil")
		log.Printf("Error: %v", err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	groupKey := constructKeyForGroup(groupName, groupVersion)
	log.Printf("Constructed group key: %s", groupKey) // Log constructed key

	pair, _, err := kv.Get(groupKey, nil)
	if err != nil {
		log.Printf("Error getting group from Consul KV: %v", err) // Log error
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	if pair == nil {
		err := fmt.Errorf("configuration group '%s' with version %.2f does not exist", groupName, groupVersion)
		log.Printf("Error: %v", err) // Log error
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	var group model.ConfigGroup
	err = json.Unmarshal(pair.Value, &group)
	if err != nil {
		log.Printf("Error unmarshalling group: %v", err) // Log error
		span.SetStatus(codes.Error, err.Error())
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
		log.Printf("Error marshalling updated group: %v", err) // Log error
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	log.Printf("Adding config to config group with SID: %s, Data: %s", groupKey, string(updatedGroupJSON)) // Log data

	p := &api.KVPair{Key: groupKey, Value: updatedGroupJSON}
	_, err = kv.Put(p, nil)
	if err != nil {
		log.Printf("Error putting updated group to Consul KV: %v", err) // Log error
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	log.Printf("Config successfully added to config group Consul KV: %s", groupKey) // Log success
	span.SetStatus(codes.Ok, "Success adding configuration group")
	return nil
}

// swagger:route DELETE /config/{name}/{groupName}/{groupVersion}/ deleteFromConfigGroup
// Delete config from group
//
// responses:
//
//	404: ErrorResponse
//	204: NoContentResponse
func (c ConfigForGroupConsulRepository) DeleteFromConfigGroup(configForGroupName string, groupName string, groupVersion float32, ctx context.Context) error {

	_, span := c.Tracer.Start(ctx, "ConfigForGroupConsulRepository.DeleteFromConfigGroup")
	defer span.End()

	kv := c.cli.KV()

	groupKey := constructKeyForGroup(groupName, groupVersion)

	pair, _, err := kv.Get(groupKey, nil)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	if pair == nil {
		span.SetStatus(codes.Error, "config does not exist")
		return fmt.Errorf("configuration group '%s' with version %.2f does not exist", groupName, groupVersion)
	}

	var group model.ConfigGroup
	err = json.Unmarshal(pair.Value, &group)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
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

		span.SetStatus(codes.Ok, "Success deleting configuration from group")
		return nil
	}

	span.SetStatus(codes.Ok, "configuration not found in the specified group")
	return errors.New("configuration not found in the specified group")
}

//func NewConfigForGroupConsulRepository() model.ConfigForGroupRepository {
//	return ConfigForGroupConsulRepository{}
//}
