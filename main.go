package main

import (
	"projekat/repositories"
	"projekat/services"
)

func main() {
	repo := repositories.NewConfigConsulRepository()
	service := services.NewConfigService(repo)
	service.Hello()
	params := make(map[string]string)
	params["username"] = "pera"
	params["password"] = "pera"
	//configs := model.NewConfig("db_config", 2, params)

	//services.AddConfig(*configs)
	//
	//handler := handlers.NewConfigHandler(services)
	//router := mux.NewRouter{}
	//router.HandleFunc("/configs/{name}/{version}", handler.Get).Methods("GET")
	//
	//http.ListenAndServe("0.0.0.0:8000", router)
	//fmt.Print("Hello")

}
