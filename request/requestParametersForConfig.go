package request

import "projekat/model"

// swagger:parameters config deleteConfig
type DeleteRequest struct {
	// Config name
	// in: path
	// required: true
	Name string `json:"group_name"`

	// Version
	// in: path
	// required: true
	Version float32 `json:"version"`
}

// swagger:parameters getConfig
type GetRequest struct {
	// Config name
	// in: path
	// required: true
	Name string `json:"group_name"`

	// Version
	// in: path
	// required: true
	Version float32 `json:"version"`
}

// swagger:parameters config addConfig
type ConfigBody struct {
	// - name: Config
	//  in: body
	//  description: name and status
	//  schema:
	//  type: object
	//     "$ref": "#/definitions/Config"
	//  required: true
	Config model.Config `json:"config"`
}
