package catalog

import (
	"encoding/json"
        "gopkg.in/yaml.v2"
        "io/ioutil"
	"net/http"
)

type Catalog struct {
	Services map[string]Service
}

type Service struct {
        DisplayName string `yaml:"display_name"`
	//Params []Param
}

func List(writer http.ResponseWriter, request *http.Request) {
        content, err := ioutil.ReadFile("catalog/catalog.yml")
        if err != nil {
		http.Error(writer, err.Error(), 500)
        }
        var catalog Catalog
        yaml.Unmarshal(content, &catalog)

	payload, err := json.Marshal(catalog)
	if err != nil {
		http.Error(writer, err.Error(), 500)
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.Write(payload)
}
