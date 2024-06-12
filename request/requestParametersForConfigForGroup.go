package request

import "projekat/model"

// swagger:parameters addToConfigGroup
type AddToConfigGroupRequest struct {
	// - name: ConfigForGroup
	//  in: body
	//  description: name and status
	//  schema:
	//  type: object
	//     "$ref": "#/definitions/ConfigForGroup"
	//  required: true
	ConfigForGroup *model.ConfigForGroup `json:"configForGroup"`
}

// swagger:parameters deleteFromConfigGroup
type DeleteFromConfigGroupRequest struct {
	// Group name
	// in: path
	// required: true
	GroupName string `json:"groupName"`

	// Group version
	// in: path
	// required: true
	GroupVersion float32 `json:"groupVersion"`
}

// swagger:parameters getConfigsByLabels
type GetConfigsByLabelsRequest struct {
	// Group name
	// in: path
	// required: true
	GroupName string `json:"groupName"`

	// Group version
	// in: path
	// required: true
	GroupVersion float32 `json:"groupVersion"`

	// Group labels
	// in: path
	// required: true
	Labels string `json:"labels"`
}

// swagger:parameters deleteConfigsByLabels
type DeleteConfigsByLabelsRequest struct {
	// Group name
	// in: path
	// required: true
	GroupName string `json:"groupName"`

	// Group version
	// in: path
	// required: true
	GroupVersion float32 `json:"groupVersion"`

	// Group labels
	// in: path
	// required: true
	Labels string `json:"labels"`
}
