package catalog

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"gopkg.in/yaml.v2"

)

// Catalog define the catalog
type Catalog struct {
	Apps map[string]App `yaml:"services"`
}

// Service define a service provided by the catalogs
type App struct {
	DisplayName string `yaml:"display_name"`
	Required    bool   `yaml:"required"`
	//Params []Param
}

var c Catalog

func Load(workdir string) {
	content, err := ioutil.ReadFile(workdir + "/catalog/catalog.yml")
	if err != nil {
		log.Print(err)
		return
	}
	yaml.Unmarshal(content, &c)

	log.Print(c)
}


// List return Services provided by the catalog
func List(writer http.ResponseWriter, request *http.Request) {

	payload, err := json.Marshal(c)
	if err != nil {
		http.Error(writer, err.Error(), 500)
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.Write(payload)
}
