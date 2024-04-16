package handlers

import (
	"projekat/services"
)

type ConfigHandler struct {
	services.ConfigService
}

//func NewConfigHandler(service services.ConfigService) ConfigHandler {
//	return ConfigHandler{
//		service: service,
//	}
//}
//
//func (c ConfigHandler) Get(w http.ResponseWriter, r *http.Request) {
//	name := mux.Vars(r)["name"]
//	version := mux.Vars(r)["version"]
//
//}
