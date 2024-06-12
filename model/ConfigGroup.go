package model

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
	GetConfigGroup(name string, version float32) (*ConfigGroup, error)
	AddConfigGroup(configGroup *ConfigGroup) error
	DeleteConfigGroup(name string, version float32) error
}
