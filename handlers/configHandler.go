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

func (c ConfigHandler) Get(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]

	versionFloat, err := strconv.ParseFloat(version, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	config, getError := c.service.GetConfig(name, float32(versionFloat))

	if getError != nil {
		http.Error(w, getError.Error(), http.StatusNotFound)
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
	versionFloat, err := strconv.ParseFloat(version, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	config, err := ch.service.GetConfig(name, float32(versionFloat))
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
