package model

type ConfigForGroup struct {
	Name       string            `json:"name"`
	Labels     map[string]string `json:"labels"`
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
	AddToConfigGroup(config *ConfigForGroup, groupName string, groupVersion float32) error
	DeleteFromConfigGroup(ConfigForGroupName string, groupName string, groupVersion float32) error
}
