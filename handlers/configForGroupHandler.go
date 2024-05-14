package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"io"
	"mime"
	"net/http"
	"projekat/model"
	"projekat/services"
	"strconv"
	"strings"
)

type ConfigForGroupHandler struct {
	service services.ConfigForGroupService
}

type AddConfigToGroupRequest struct {
	ConfigForGroup struct {
		Name       string            `json:"name"`
		Labels     map[string]string `json:"labels"`
		Parameters map[string]string `json:"parameters"`
	} `json:"config"`
	ConfigGroup struct {
		Name           string                 `json:"name"`
		Version        float32                `json:"version"`
		Configurations []model.ConfigForGroup `json:"configurations"`
	} `json:"configGroup"`
}

func NewConfigForGroupHandler(service services.ConfigForGroupService) ConfigForGroupHandler {
	return ConfigForGroupHandler{
		service: service,
	}
}

func decoder(ctx context.Context, r io.Reader) (*model.ConfigForGroup, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var c model.ConfigForGroup
	if err := dec.Decode(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

func renderer(ctx context.Context, w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		return
	}
}

func (ch *ConfigForGroupHandler) AddToConfigGroup(w http.ResponseWriter, req *http.Request) {

	contentType := req.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, "dipshit: "+err.Error(), http.StatusBadRequest)
		return
	}
	if mediaType != "application/json" {
		err := errors.New("expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	var addToGroupReq AddConfigToGroupRequest
	err = json.NewDecoder(req.Body).Decode(&addToGroupReq)
	if err != nil {
		http.Error(w, "failed to decode JSON request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = ch.service.AddToConfigGroup(addToGroupReq.ConfigForGroup.Name, addToGroupReq.ConfigForGroup.Labels, addToGroupReq.ConfigForGroup.Parameters, addToGroupReq.ConfigGroup.Name, addToGroupReq.ConfigGroup.Version)
	if err != nil {
		http.Error(w, "Failed to add configuration to configuration group: "+err.Error(), http.StatusInternalServerError)
		return
	}

	renderer(req.Context(), w, map[string]string{"message": "Configuration added to group successfully"})
}

func (ch *ConfigForGroupHandler) DeleteFromConfigGroup(w http.ResponseWriter, req *http.Request) {
	configForGroupName := mux.Vars(req)["name"]
	groupName := mux.Vars(req)["groupName"]
	groupVersion := mux.Vars(req)["groupVersion"]

	versionFloat1, err := strconv.ParseFloat(groupVersion, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	groupVersion32 := float32(versionFloat1)

	err = ch.service.DeleteFromConfigGroup(configForGroupName, groupName, groupVersion32)
	if err != nil {
		http.Error(w, "Failed to delete configuration from configuration group: "+err.Error(), http.StatusInternalServerError)
		return
	}

	renderer(req.Context(), w, map[string]string{"message": "Configuration deleted from group successfully"})

}

func (ch *ConfigForGroupHandler) GetConfigsByLabels(w http.ResponseWriter, req *http.Request) {
	groupName := mux.Vars(req)["groupName"]
	groupVersion := mux.Vars(req)["groupVersion"]
	labels := mux.Vars(req)["labels"]

	versionFloat1, err := strconv.ParseFloat(groupVersion, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	groupVersion32 := float32(versionFloat1)

	labelMap := make(map[string]string)
	for _, label := range strings.Split(labels, ",") {
		parts := strings.SplitN(label, ":", 2)
		if len(parts) == 2 {
			labelMap[parts[0]] = parts[1]
		}
	}

	// Call the service method to get configurations by labels
	configs, err := ch.service.GetConfigsByLabels(groupName, groupVersion32, labelMap)
	if err != nil {
		http.Error(w, "Failed to get configurations by labels from configuration group: "+err.Error(), http.StatusInternalServerError)
		return
	}

	configMap := make(map[string]interface{})
	for _, config := range configs {
		configMap[config.Name] = config
	}

	renderer(req.Context(), w, configMap)

}
