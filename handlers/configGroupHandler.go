package handlers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"mime"
	"net/http"
	"projekat/model"
	"projekat/services"
	"strconv"
	"strings"
)

type ConfigGroupHandler struct {
	service services.ConfigGroupService
}

func NewConfigGroupHandler(service services.ConfigGroupService) ConfigGroupHandler {
	return ConfigGroupHandler{
		service: service,
	}
}

func (ch *ConfigGroupHandler) CreateConfigGroup(w http.ResponseWriter, req *http.Request) {
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
	var configGroup model.ConfigGroup
	err = json.NewDecoder(req.Body).Decode(&configGroup)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ch.service.AddConfigGroup(configGroup.Name, configGroup.Version, configGroup.Configurations)
	if err != nil {
		// Log the error for debugging purposes
		log.Printf("Error adding config group: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderJSON(req.Context(), w, configGroup)
}

func (c *ConfigGroupHandler) GetConfigGroup(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]

	versionFloat, err := strconv.ParseFloat(version, 64) // ParseFloat returns float64
	if err != nil {
		http.Error(w, "Invalid version number", http.StatusBadRequest)
		return
	}
	// Convert float64 to float32
	version32 := float32(versionFloat)

	config, err := c.service.GetConfigGroup(name, version32)
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
