package request

import "projekat/model"

// swagger:response ResponseConfig
type ResponseConfig struct {
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

// swagger:response ResponseConfigGroup
type ResponseConfigGroup struct {
	// Name of the ConfigGroup
	// in: string
	Name string `json:"name"`

	// Version of the ConfigGroup
	// in: float32
	Version float32 `json:"version"`

	// Configurations of the ConfigGroup
	// in: []ConfigForGroup
	Configurations []model.ConfigForGroup `json:"configurations"`
}

// swagger:response ResponseConfigForGroup
type ResponseConfigForGroup struct {
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

// swagger:response ErrorResponse
type ErrorResponse struct {
	// Error status code
	// in: int64
	Status int64 `json:"status"`
	// Message of the error
	// in: string
	Message string `json:"message"`
}

// swagger:response NoContentResponse
type NoContentResponse struct{}
