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

type ConfigConsulRepository struct {
	cli    *api.Client
	logger *log.Logger
	Tracer trace.Tracer
}

func New(logger *log.Logger, tracer trace.Tracer) (*ConfigConsulRepository, error) {
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
	return &ConfigConsulRepository{cli: client, logger: logger, Tracer: tracer}, nil
}

// swagger:route GET /config/{name}/{version}/ getConfig
// Get config
//
// responses:
//
//	200: ResponseConfig
func (c ConfigConsulRepository) GetConfig(name string, version float32, ctx context.Context) (*model.Config, error) {
	_, span := c.Tracer.Start(ctx, "ConfigConsulRepository.GetConfig")
	defer span.End()

	if c.cli == nil {
		err := errors.New("consul client is nil")
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

	key := constructKey(name, version)
	log.Printf("Constructed key: %s", key) // Log constructed key

	pair, _, err := kv.Get(key, nil)
	if err != nil {
		log.Printf("Error getting key from KV store: %v", err) // Log error
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	if pair == nil {
		err := fmt.Errorf("configuration '%s' with version %.1f not found", name, version)
		log.Printf("Error: %v", err) // Log error
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	log.Printf("KV pair retrieved: %s", pair.Value) // Log retrieved KV pair

	config := &model.Config{}
	err = json.Unmarshal(pair.Value, config)
	if err != nil {
		log.Printf("Error unmarshalling KV pair value: %v", err) // Log error
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	log.Printf("Config unmarshalled: %+v", config) // Log unmarshalled config

	span.SetStatus(codes.Ok, "Success getting configuration")
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
func (c ConfigConsulRepository) AddConfig(config *model.Config, ctx context.Context) error {
	_, span := c.Tracer.Start(ctx, "ConfigConsulRepository.AddConfig")
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

	key := constructKey(config.Name, config.Version)
	log.Printf("Constructed key: %s", key) // Log constructed key

	data, err := json.Marshal(config)
	if err != nil {
		log.Printf("Error marshalling config: %v", err) // Log error
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	log.Printf("Adding config with SID: %s, Data: %s", key, string(data)) // Log data

	p := &api.KVPair{Key: key, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		log.Printf("Error putting config to Consul KV: %v", err) // Log error
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	log.Printf("Config successfully added to Consul KV: %s", key) // Log success
	span.SetStatus(codes.Ok, "Config successfully added")
	return nil
}

// swagger:route DELETE /config/{name}/{version}/ config deleteConfig
// Delete config
//
// responses:
//
//	404: ErrorResponse
//	204: NoContentResponse
func (c ConfigConsulRepository) DeleteConfig(name string, version float32, ctx context.Context) error {
	_, span := c.Tracer.Start(ctx, "ConfigRepository.DeleteConfig")
	defer span.End()

	kv := c.cli.KV()
	_, err := kv.Delete(constructKey(name, version), nil)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	c.logger.Println("Config successfully deleted from Consul:", name)

	span.SetStatus(codes.Ok, "Successfully deleted configuration group")
	return nil
}

func (cr *ConfigConsulRepository) GetIdempotencyRequestByKey(key string, ctx context.Context) (bool, error) {
	_, span := cr.Tracer.Start(ctx, "Repository.GetIdempotencyRequest")
	defer span.End()
	kv := cr.cli.KV()

	data, _, err := kv.Get(constructIdempotencyRequestKey(key), nil)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return false, err
	}

	if data == nil {
		span.SetStatus(codes.Error, "Idempotency request not found") // Set appropriate status message
		return false, nil
	}

	span.SetStatus(codes.Ok, "Success")
	return true, nil
}

func (cr *ConfigConsulRepository) AddIdempotencyRequest(req *model.IdempotencyRequest, ctx context.Context) (*model.IdempotencyRequest, error) {
	_, span := cr.Tracer.Start(ctx, "Repository.AddIdempotencyRequest")

	kv := cr.cli.KV()

	data, err := json.Marshal(req)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	keyValue := &api.KVPair{Key: constructIdempotencyRequestKey(req.Key), Value: data}
	_, err = kv.Put(keyValue, nil)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "Success")
	return req, nil
}

//func NewConfigConsulRepository() model.ConfigRepository {
//	return ConfigConsulRepository{}
//}
