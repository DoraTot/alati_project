package model

type Config struct {
	// todo: dodati atribute
	Name       string
	Version    float32
	Parameters map[string]string
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
}
