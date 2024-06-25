package repositories

import (
	"context"
	"encoding/json"
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

	kv := c.cli.KV()
	key := constructKey(name, version)
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return nil, err
	}
	if pair == nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, fmt.Errorf("configuration '%s' with version %.1f not found", name, version)
	}
	config := &model.Config{}
	err = json.Unmarshal(pair.Value, config)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

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
	kv := c.cli.KV()
	key := constructKey(config.Name, config.Version)

	data, err := json.Marshal(config)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	c.logger.Printf("Adding config with SID: %s, Data: %s\n", key, string(data))

	p := &api.KVPair{Key: key, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		c.logger.Println("Error putting config to Consul KV:", err)
		return err
	}
	c.logger.Println("Config successfully added to Consul KV:", key)
	span.SetStatus(codes.Ok, "Success getting configuration")
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
		span.SetStatus(codes.Error, err.Error())
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
