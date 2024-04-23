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

type ConfigHandler struct {
	service services.ConfigService
}

func NewConfigHandler(service services.ConfigService) ConfigHandler {
	return ConfigHandler{
		service: service,
	}
}

func decodeBody(ctx context.Context, r io.Reader) (*model.Config, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var c model.Config
	if err := dec.Decode(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

func renderJSON(ctx context.Context, w http.ResponseWriter, v interface{}) {
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

func (c *ConfigHandler) Get(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]

	versionFloat, err := strconv.ParseFloat(version, 64) // ParseFloat returns float64
	if err != nil {
		http.Error(w, "Invalid version number", http.StatusBadRequest)
		return
	}
	// Convert float64 to float32
	version32 := float32(versionFloat)

	config, err := c.service.GetConfig(name, version32)
	if err != nil {
		if strings.Contains(err.Error(), "config not found") {
			http.Error(w, "Configuration not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve configuration", http.StatusInternalServerError)
		}
		return
	}

	renderJSON(r.Context(), w, config)

}

func (ch *ConfigHandler) CreatePostHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediaType != "application/json" {
		err := errors.New("expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}
	config, err := decodeBody(req.Context(), req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = ch.service.AddConfig(config.Name, config.Version, config.Parameters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderJSON(req.Context(), w, config)
}

func (ch *ConfigHandler) DelPostHandler(w http.ResponseWriter, req *http.Request) {
	name := mux.Vars(req)["name"]
	version := mux.Vars(req)["version"]
	versionFloat, err := strconv.ParseFloat(version, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	version32 := float32(versionFloat)

	config, err := ch.service.GetConfig(name, version32)
	if err != nil {
		http.Error(w, "Configuration not found: "+err.Error(), http.StatusNotFound)
		return
	}

	err = ch.service.DeleteConfig(config.Name, config.Version)
	if err != nil {
		http.Error(w, "Failed to delete configuration: "+err.Error(), http.StatusInternalServerError)
		return
	}

	renderJSON(req.Context(), w, map[string]string{"message": "Configuration deleted successfully"})

}

func (ch *ConfigHandler) AddToConfigGroup(w http.ResponseWriter, req *http.Request) {
	name := mux.Vars(req)["name"]
	version := mux.Vars(req)["version"]
	groupName := mux.Vars(req)["groupName"]
	groupVersion := mux.Vars(req)["groupVersion"]

	versionFloat, err := strconv.ParseFloat(version, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	version32 := float32(versionFloat)

	versionFloat1, err := strconv.ParseFloat(groupVersion, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	groupVersion32 := float32(versionFloat1)

	config, err := ch.service.GetConfig(name, version32)
	if err != nil {
		http.Error(w, "Configuration not found: "+err.Error(), http.StatusNotFound)
		return
	}

	err = ch.service.AddToConfigGroup(config, groupName, groupVersion32)
	if err != nil {
		http.Error(w, "Failed to add configuration to configuration group: "+err.Error(), http.StatusInternalServerError)
		return
	}

	renderJSON(req.Context(), w, map[string]string{"message": "Configuration added to group successfully"})

}

func (ch *ConfigHandler) DeleteFromConfigGroup(w http.ResponseWriter, req *http.Request) {
	name := mux.Vars(req)["name"]
	version := mux.Vars(req)["version"]
	groupName := mux.Vars(req)["groupName"]
	groupVersion := mux.Vars(req)["groupVersion"]

	versionFloat, err := strconv.ParseFloat(version, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	version32 := float32(versionFloat)

	versionFloat1, err := strconv.ParseFloat(groupVersion, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	groupVersion32 := float32(versionFloat1)

	config, err := ch.service.GetConfig(name, version32)
	if err != nil {
		http.Error(w, "Configuration not found: "+err.Error(), http.StatusNotFound)
		return
	}

	err = ch.service.DeleteFromConfigGroup(config, groupName, groupVersion32)
	if err != nil {
		http.Error(w, "Failed to delete configuration from configuration group: "+err.Error(), http.StatusInternalServerError)
		return
	}

	renderJSON(req.Context(), w, map[string]string{"message": "Configuration deleted from group successfully"})

}
