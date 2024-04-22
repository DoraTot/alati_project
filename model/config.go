package model

type Config struct {
	// todo: dodati atribute
	Name       string            `json:"name"`
	Version    float32           `json:"version"`
	Parameters map[string]string `json:"parameters"`
}

// todo: dodati metode

func NewConfig(name string, version float32, parameters map[string]string) *Config {
	return &Config{
		Name:       name,
		Version:    version,
		Parameters: parameters,
	}
}

type ConfigRepository interface {
	// todo: dodati metode

	GetConfig(name string, version float32) (*Config, error)
	AddConfig(config *Config) error
	DeleteConfig(name string) error
	AddToConfigGroup(config *Config, groupName string) error
	DeleteFromConfigGroup(config *Config, groupName string) error
}
