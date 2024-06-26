package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"log"
	"mime"
	"net/http"
	"projekat/model"
	"projekat/services"
	"strconv"
	"strings"
)

type ConfigGroupHandler struct {
	Service services.ConfigGroupService
	Tracer  trace.Tracer
}

func NewConfigGroupHandler(service services.ConfigGroupService, tracer trace.Tracer) ConfigGroupHandler {
	return ConfigGroupHandler{
		service,
		tracer,
	}
}

func (ch *ConfigGroupHandler) CreateConfigGroup(w http.ResponseWriter, req *http.Request) {
	log.Println("Entered CreateConfigGroup")
	ctx, span := ch.Tracer.Start(req.Context(), "ConfigGroupHandler.CreateConfigGroup")
	defer span.End()

	contentType := req.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediaType != "application/json" {
		span.SetStatus(codes.Error, err.Error())
		err := errors.New("expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}
	var configGroup model.ConfigGroup
	err = json.NewDecoder(req.Body).Decode(&configGroup)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ch.Service.AddConfigGroup(configGroup.Name, configGroup.Version, configGroup.Configurations, ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		// Log the error for debugging purposes
		log.Printf("Error adding config group: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//renderJSON(req.Context(), w, configGroup)
	renderJSON(ctx, w, configGroup)
	span.SetStatus(codes.Ok, "")
}

func (c *ConfigGroupHandler) GetConfigGroup(w http.ResponseWriter, r *http.Request) {

	ctx, span := c.Tracer.Start(r.Context(), "ConfigGroupHandler.GetConfigGroup")
	defer span.End()

	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]

	versionFloat, err := strconv.ParseFloat(version, 64) // ParseFloat returns float64
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, "Invalid version number", http.StatusBadRequest)
		return
	}
	// Convert float64 to float32
	version32 := float32(versionFloat)

	config, err := c.Service.GetConfigGroup(name, version32, ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		errMsg := fmt.Sprintf("configGroup '%s' with version %.2f not found", name, version)
		if strings.Contains(err.Error(), errMsg) {
			http.Error(w, "Configuration group not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve configuration group", http.StatusInternalServerError)
		}
		return
	}

	//renderJSON(r.Context(), w, config)
	renderJSON(ctx, w, config)
	span.SetStatus(codes.Ok, "")

}

func (ch *ConfigGroupHandler) DeleteConfigGroup(w http.ResponseWriter, req *http.Request) {

	ctx, span := ch.Tracer.Start(req.Context(), "ConfigGroupHandler.DeleteConfigGroup")
	defer span.End()

	name := mux.Vars(req)["name"]
	version := mux.Vars(req)["version"]
	versionFloat, err := strconv.ParseFloat(version, 64)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	version32 := float32(versionFloat)

	configGroup, err := ch.Service.GetConfigGroup(name, version32, ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, "Configuration group not found: "+err.Error(), http.StatusNotFound)
		return
	}

	err = ch.Service.DeleteConfigGroup(configGroup.Name, configGroup.Version, ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, "Failed to delete configuration group: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//renderJSON(req.Context(), w, map[string]string{"message": "Configuration group deleted successfully"})
	renderJSON(ctx, w, map[string]string{"message": "Configuration group deleted successfully"})
	span.SetStatus(codes.Ok, "")
}
