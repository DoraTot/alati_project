package model

// swagger:model Config
type Config struct {
	// Name of the Config
	// in: string
	Name string `json:"name"`

	// Version of the Config
	// in: float32
	Version float32 `json:"version"`

	// Parameters of the Config
	// in: map[string]string
	Parameters map[string]string `json:"parameters"`
}

func NewConfig(name string, version float32, parameters map[string]string) *Config {
	return &Config{
		Name:       name,
		Version:    version,
		Parameters: parameters,
	}
}

type ConfigRepository interface {
	GetConfig(name string, version float32) (*Config, error)
	AddConfig(config *Config) error
	DeleteConfig(name string, version float32) error
}
