package model

type ConfigGroup struct {
	Name           string           `json:"name"`
	Version        float32          `json:"version"`
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
	GetConfigGroup(name string, version float32) (*ConfigGroup, error)
	AddConfigGroup(configGroup *ConfigGroup) error
	DeleteConfigGroup(name string, version float32) error
}
