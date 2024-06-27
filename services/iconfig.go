// services/config_service_interface.go
package services

import (
	"context"
	"projekat/model"
)

type ConfigServiceInterface interface {
	Hello()
	AddConfig(name string, version float32, parameters map[string]string, ctx context.Context) error
	GetConfig(name string, version float32, ctx context.Context) (*model.Config, error)
	DeleteConfig(name string, version float32, ctx context.Context) error
}
