package model

type ConfigForGroup struct {
	Labels     map[string]string `json:"labels"`
	Name       string            `json:"name"`
	Parameters map[string]string `json:"parameters"`
}

// todo: dodati metode

func NewConfigForGroup(name string, labels map[string]string, parameters map[string]string) *ConfigForGroup {
	return &ConfigForGroup{
		Name:       name,
		Labels:     labels,
		Parameters: parameters,
	}
}

type ConfigForGroupRepository interface {
	GetConfig(name string, version float32) (*ConfigForGroup, error)
	AddConfig(config *ConfigForGroup) error
	DeleteConfig(name string, version float32) error
	AddToConfigGroup(config *ConfigForGroup, groupName string, groupVersion float32) error
	DeleteFromConfigGroup(config *ConfigForGroup, groupName string, groupVersion float32) error
}
