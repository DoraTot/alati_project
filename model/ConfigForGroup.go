package model

import "context"

// swagger:model ConfigForGroup
type ConfigForGroup struct {
	// Name of the ConfigForGroup
	// in: string
	Name string `json:"name"`

	// Labels of the ConfigForGroup
	// in: map[string]string
	Labels map[string]string `json:"labels"`

	// Parameters of the ConfigForGroup
	// in: map[string]string
	Parameters map[string]string `json:"parameters"`
}

func NewConfigForGroup(name string, labels map[string]string, parameters map[string]string) *ConfigForGroup {
	return &ConfigForGroup{
		Name:       name,
		Labels:     labels,
		Parameters: parameters,
	}
}

type ConfigForGroupRepository interface {
	AddToConfigGroup(config *ConfigForGroup, groupName string, groupVersion float32, ctx context.Context) error
	DeleteFromConfigGroup(ConfigForGroupName string, groupName string, groupVersion float32, ctx context.Context) error
	GetConfigsByLabels(groupName string, groupVersion float32, labels map[string]string, ctx context.Context) ([]ConfigForGroup, error)
	DeleteConfigsByLabels(groupName string, groupVersion float32, labels map[string]string, ctx context.Context) error
}
