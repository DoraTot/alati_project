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

type ConfigGroupConsulRepository struct {
	cli    *api.Client
	logger *log.Logger
	Tracer trace.Tracer
}

func NewCG(logger *log.Logger, trace trace.Tracer) (*ConfigGroupConsulRepository, error) {
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
	return &ConfigGroupConsulRepository{cli: client, logger: logger, Tracer: trace}, nil
}

// swagger:route GET /configGroup/{name}/{version}/ getConfigGroup
// Get configGroup by ID
//
// responses:
//
//	404: ErrorResponse
//	200: ResponseConfigGroup
func (c ConfigGroupConsulRepository) GetConfigGroup(name string, version float32, ctx context.Context) (*model.ConfigGroup, error) {
	_, span := c.Tracer.Start(ctx, "ConfigGroupConsulRepository.GetConfigGroup")
	defer span.End()

	if c.cli == nil {
		err := errors.New("Consul client is nil")
		log.Printf("Error: %v", err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	kv := c.cli.KV()
	if kv == nil {
		err := errors.New("KV store is nil")
		log.Printf("Error: %v", err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	key := constructKeyForGroup(name, version)
	log.Printf("Constructed group key: %s", key) // Log constructed key

	pair, _, err := kv.Get(key, nil)
	if err != nil {
		log.Printf("Error getting config group from Consul KV: %v", err) // Log error
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	if pair == nil {
		err := fmt.Errorf("configuration group '%s' with version %.1f not found", name, version)
		log.Printf("Error: %v", err) // Log error
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	configGroup := &model.ConfigGroup{}
	err = json.Unmarshal(pair.Value, configGroup)
	if err != nil {
		log.Printf("Error unmarshalling config group: %v", err) // Log error
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	log.Printf("Retrieved config group: %+v", configGroup) // Log retrieved config group
	span.SetStatus(codes.Ok, "Success getting config group")
	return configGroup, nil
}

// swagger:route POST /configGroup/ configGroup addConfigGroup
// Add new configGroup
//
// responses:
//
//	415: ErrorResponse
//	400: ErrorResponse
//	201: ResponseConfigGroup
func (c ConfigGroupConsulRepository) AddConfigGroup(config *model.ConfigGroup, ctx context.Context) error {
	_, span := c.Tracer.Start(ctx, "ConfigGroupConsulRepository.AddConfigGroup")
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

	key := constructKeyForGroup(config.Name, config.Version)
	log.Printf("Constructed group key: %s", key) // Log constructed key

	data, err := json.Marshal(config)
	if err != nil {
		log.Printf("Error marshalling config group: %v", err) // Log error
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	log.Printf("Adding config group with SID: %s, Data: %s", key, string(data)) // Log data being added

	p := &api.KVPair{Key: key, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		log.Printf("Error adding config group to Consul KV: %v", err) // Log error
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	log.Printf("Config group added successfully: %s", key) // Log success
	span.SetStatus(codes.Ok, "Config group added successfully")
	return nil
}

// swagger:route DELETE /configGroup/{name}/{version}/ deleteConfigGroup
// Delete configGroup
//
// responses:
//
//	404: ErrorResponse
//	204: NoContentResponse
func (c ConfigGroupConsulRepository) DeleteConfigGroup(name string, version float32, ctx context.Context) error {
	_, span := c.Tracer.Start(ctx, "ConfigGroupConsulRepository.DeleteConfigGroup")
	defer span.End()
	kv := c.cli.KV()
	_, err := kv.Delete(constructKeyForGroup(name, version), nil)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		c.logger.Println("Error deleting config group:", err)
		return err
	}

	c.logger.Println("Config group deleted successfully", constructKeyForGroup(name, version))
	span.SetStatus(codes.Ok, "Config group deleted successfully")
	return nil
}

//func NewConfigGroupConsulRepository() model.ConfigGroupRepository {
//	return ConfigGroupConsulRepository{}
//}
