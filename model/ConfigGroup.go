package model

import "context"

// swagger:model ConfigGroup
type ConfigGroup struct {
	// Name of the ConfigGroup
	// in: string
	Name string `json:"name"`

	// Version of the ConfigGroup
	// in: float32
	Version float32 `json:"version"`

	// Configurations of the ConfigGroup
	// in: []ConfigForGroup
	Configurations []ConfigForGroup `json:"configurations"`
}

func NewConfigGroup(name string, version float32, configurations []ConfigForGroup) *ConfigGroup {
	return &ConfigGroup{
		Name:           name,
		Version:        version,
		Configurations: configurations,
	}
}

type ConfigGroupRepository interface {
	GetConfigGroup(name string, version float32, ctx context.Context) (*ConfigGroup, error)
	AddConfigGroup(configGroup *ConfigGroup, ctx context.Context) error
	DeleteConfigGroup(name string, version float32, ctx context.Context) error
}
