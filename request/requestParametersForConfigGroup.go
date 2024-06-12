package request

import "projekat/model"

// swagger:parameters deleteConfigGroup
type DeleteGroupRequest struct {
	// Group name
	// in: path
	// required: true
	Name string `json:"group_name"`

	// Version
	// in: path
	// required: true
	Version float32 `json:"version"`
}

// swagger:parameters getConfigGroup
type GetGroupRequest struct {
	// Group name
	// in: path
	// required: true
	Name string `json:"group_name"`

	// Version
	// in: path
	// required: true
	Version float32 `json:"version"`
}

// swagger:parameters getConfigGroup
type ConfigGroupBody struct {
	// - name: ConfigGroup
	//  in: body
	//  description: config group
	//  schema:
	//  type: object
	//     "$ref": "#/definitions/ConfigGroup"
	//  required: true
	ConfigGroup model.ConfigGroup `json:"configGroup"`
}
