package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"io"
	"log"
	"mime"
	"net/http"
	"projekat/model"
	"projekat/services"
	"strconv"
	"strings"
)

type ConfigHandler struct {
	logger  *log.Logger
	Service services.ConfigService
	Tracer  trace.Tracer
}

func NewConfigHandler(l *log.Logger, s services.ConfigService, tracer trace.Tracer) *ConfigHandler {
	return &ConfigHandler{l, s, tracer}
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
		log.Println("There has been an internal error.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		// Log the error instead of returning it
		log.Println("Error writing response:", err)
	}
}

func (c *ConfigHandler) Get(w http.ResponseWriter, r *http.Request) {
	log.Println("Entering Get handler") // Log entry

	ctx, span := c.Tracer.Start(r.Context(), "ConfigHandler.Get")
	defer span.End()

	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]
	log.Printf("Received request for config: name=%s, version=%s", name, version) // Log request details

	versionFloat, err := strconv.ParseFloat(version, 64) // ParseFloat returns float64
	if err != nil {
		log.Printf("Error parsing version: %v", err) // Log error
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, "Invalid version number", http.StatusBadRequest)
		return
	}
	log.Printf("Parsed version: %f", versionFloat) // Log parsed version

	// Convert float64 to float32
	version32 := float32(versionFloat)
	log.Printf("Converted version to float32: %f", version32) // Log converted version

	config, err := c.Service.GetConfig(name, version32, ctx)
	if err != nil {
		log.Printf("Error getting config: %v", err) // Log error
		span.SetStatus(codes.Error, err.Error())
		if strings.Contains(err.Error(), "config not found") {
			http.Error(w, "Configuration not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve configuration", http.StatusInternalServerError)
		}
		return
	}
	log.Printf("Retrieved config: %+v", config) // Log retrieved config

	renderJSON(ctx, w, config)
	span.SetStatus(codes.Ok, "")
	log.Println("Successfully processed Get handler") // Log successful completion
}

func (ch *ConfigHandler) CreatePostHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("Entering Post handler")
	ctx, span := ch.Tracer.Start(req.Context(), "ConfigHandler.CreatePostHandler")
	defer span.End()

	// Retrieve Content-Type header and validate media type
	contentType := req.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		log.Printf("Error parsing Content-Type header: %v", err)
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediaType != "application/json" {
		err := errors.New("expect application/json Content-Type")
		log.Printf("Invalid media type: %s", mediaType)
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	// Decode request body into model.Config
	config, err := decodeBody(req.Context(), req.Body)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Decoded config: %+v", config)

	// Call service to add configuration
	err = ch.Service.AddConfig(config.Name, config.Version, config.Parameters, ctx)
	if err != nil {
		log.Printf("Error adding config: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render response as JSON
	renderJSON(ctx, w, config)
	log.Println("Successfully processed CreatePostHandler")
	span.SetStatus(codes.Ok, "")
}

func (ch *ConfigHandler) DelPostHandler(w http.ResponseWriter, req *http.Request) {
	ctx, span := ch.Tracer.Start(req.Context(), "ConfigHandler.DelPostHandler")
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

	config, err := ch.Service.GetConfig(name, version32, ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, "Configuration not found: "+err.Error(), http.StatusNotFound)
		return
	}

	err = ch.Service.DeleteConfig(config.Name, config.Version, ctx)
	if err != nil {
		http.Error(w, "Failed to delete configuration: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//renderJSON(req.Context(), w, map[string]string{"message": "Configuration deleted successfully"})
	renderJSON(ctx, w, map[string]string{"message": "Configuration deleted successfully"})
	span.SetStatus(codes.Ok, "")
}
