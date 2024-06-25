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

	ctx, span := c.Tracer.Start(r.Context(), "ConfigHandler.Get")
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

	config, err := c.Service.GetConfig(name, version32, ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		if strings.Contains(err.Error(), "config not found") {
			http.Error(w, "Configuration not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve configuration", http.StatusInternalServerError)
		}
		return
	}

	renderJSON(ctx, w, config)
	span.SetStatus(codes.Ok, "")

}

func (ch *ConfigHandler) CreatePostHandler(w http.ResponseWriter, req *http.Request) {
	ctx, span := ch.Tracer.Start(req.Context(), "ConfigHandler.CreateConfiguration")
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
	config, err := decodeBody(req.Context(), req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = ch.Service.AddConfig(config.Name, config.Version, config.Parameters, ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderJSON(ctx, w, config)
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
