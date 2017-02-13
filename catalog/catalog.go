package catalog

import (
	"encoding/json"
	"net/http"
)

type Catalog struct {
	Services []Service
}

type Service struct {
	Name string
	//Params []Param
}

func List(writer http.ResponseWriter, request *http.Request) {
	app := Service{"test"}
	payload, err := json.Marshal(app)
	if err != nil {
		http.Error(writer, err.Error(), 500)
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.Write(payload)
}
