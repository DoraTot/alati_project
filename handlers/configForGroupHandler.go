package handlers

import (
	"bytes"
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

type ConfigForGroupHandler struct {
	Service services.ConfigForGroupService
	Tracer  trace.Tracer
}

type AddConfigToGroupRequest struct {
	ConfigForGroup struct {
		Name       string            `json:"name"`
		Labels     map[string]string `json:"labels"`
		Parameters map[string]string `json:"parameters"`
	}
	ConfigGroup struct {
		Name           string                 `json:"name"`
		Version        float32                `json:"version"`
		Configurations []model.ConfigForGroup `json:"configurations"`
	} `json:"configGroup"`
}

func NewConfigForGroupHandler(service services.ConfigForGroupService, tracer trace.Tracer) ConfigForGroupHandler {
	return ConfigForGroupHandler{
		service,
		tracer,
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

func (ch *ConfigForGroupHandler) renderer(ctx context.Context, w http.ResponseWriter, v interface{}) {
	ctx, span := ch.Tracer.Start(ctx, "renderJSON")
	defer span.End()

	js, err := json.Marshal(v)
	if err != nil {

		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if _, err := w.Write(js); err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (ch *ConfigForGroupHandler) AddToConfigGroup(w http.ResponseWriter, req *http.Request) {
	ctx, span := ch.Tracer.Start(req.Context(), "ConfigForGroupHandler.AddToConfigGroup")
	defer span.End()

	vars := mux.Vars(req)
	groupName := vars["groupName"]
	groupVersionStr := vars["groupVersion"]

	log.Printf("The version for configGroup is %s", groupVersionStr)

	groupVersion, err := strconv.ParseFloat(groupVersionStr, 32)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, "invalid groupVersion", http.StatusBadRequest)
		return
	}

	contentType := req.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, "an error has occurred: "+err.Error(), http.StatusBadRequest)
		return
	}
	if mediaType != "application/json" {
		span.SetStatus(codes.Error, err.Error())
		err := errors.New("expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	// Read the entire request body
	body, err := io.ReadAll(req.Body)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, "failed to read request body: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Add the wrapper around the JSON payload
	wrappedBody := []byte(`{"configForGroup":` + string(body) + `}`)

	// Create a new request with the wrapped body
	req.Body = io.NopCloser(bytes.NewReader(wrappedBody))

	var addToGroupReq AddConfigToGroupRequest
	err = json.NewDecoder(req.Body).Decode(&addToGroupReq)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, "failed to decode JSON request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Assuming addToGroupReq.ConfigForGroup is of type model.ConfigForGroup
	err = ch.Service.AddToConfigGroup(addToGroupReq.ConfigForGroup.Name, addToGroupReq.ConfigForGroup.Labels, addToGroupReq.ConfigForGroup.Parameters, groupName, float32(groupVersion), ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, "Failed to add configuration to configuration group: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ch.renderer(ctx, w, map[string]string{"message": "Configuration added to group successfully"})
	span.SetStatus(codes.Ok, "")
}

func (ch *ConfigForGroupHandler) DeleteFromConfigGroup(w http.ResponseWriter, req *http.Request) {
	ctx, span := ch.Tracer.Start(req.Context(), "ConfigForGroupHandler.DeleteFromConfigGroup")
	defer span.End()

	configForGroupName := mux.Vars(req)["name"]
	groupName := mux.Vars(req)["groupName"]
	groupVersion := mux.Vars(req)["groupVersion"]

	versionFloat1, err := strconv.ParseFloat(groupVersion, 64)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	groupVersion32 := float32(versionFloat1)

	err = ch.Service.DeleteFromConfigGroup(configForGroupName, groupName, groupVersion32, ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, "Failed to delete configuration from configuration group: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ch.renderer(ctx, w, map[string]string{"message": "Configuration deleted from group successfully"})
	span.SetStatus(codes.Ok, "")
}

func (ch *ConfigForGroupHandler) GetConfigsByLabels(w http.ResponseWriter, req *http.Request) {
	ctx, span := ch.Tracer.Start(req.Context(), "ConfigForGroupHandler.GetConfigsByLabels")
	defer span.End()

	vars := mux.Vars(req)
	groupName := vars["groupName"]
	groupVersionStr := vars["groupVersion"]
	labelsStr := vars["labels"]

	log.Printf("groupName: %s, groupVersion: %s, labels: %s", groupName, groupVersionStr, labelsStr)

	// Parse groupVersion from string to float32
	groupVersion, err := strconv.ParseFloat(groupVersionStr, 32)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, "invalid groupVersion", http.StatusBadRequest)
		return
	}
	groupVersion32 := float32(groupVersion)

	// Parse labels from string to map[string]string
	labelMap := make(map[string]string)
	for _, label := range strings.Split(labelsStr, ",") {
		parts := strings.SplitN(label, ":", 2)
		if len(parts) == 2 {
			labelMap[parts[0]] = parts[1]
		}
	}

	// Call the service method to get configurations by labels
	configs, err := ch.Service.GetConfigsByLabels(groupName, groupVersion32, labelMap, ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, "Failed to get configurations by labels from configuration group: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare response
	configMap := make(map[string]interface{})
	for _, config := range configs {
		configMap[config.Name] = config
	}

	ch.renderer(ctx, w, configMap)
	span.SetStatus(codes.Ok, "")
}

func (ch *ConfigForGroupHandler) DeleteConfigsByLabels(w http.ResponseWriter, req *http.Request) {
	ctx, span := ch.Tracer.Start(req.Context(), "ConfigForGroupHandler.DeleteConfigByLabels")
	defer span.End()

	groupName := mux.Vars(req)["groupName"]
	groupVersion := mux.Vars(req)["groupVersion"]
	labels := mux.Vars(req)["labels"]

	versionFloat1, err := strconv.ParseFloat(groupVersion, 64)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
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

	err = ch.Service.DeleteConfigsByLabels(groupName, groupVersion32, labelMap, ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, "Failed to delete configuration from configuration group: "+err.Error(), http.StatusInternalServerError)
		return
	}
	ch.renderer(ctx, w, map[string]string{"message": "Configuration deleted from group successfully"})
	span.SetStatus(codes.Ok, "")
}
