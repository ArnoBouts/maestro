package catalog

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"gopkg.in/yaml.v2"

	"golang.org/x/net/context"

	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/docker/ctx"
	"github.com/docker/libcompose/project"

	"github.com/gorilla/mux"
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

func Load() Catalog {
	content, err := ioutil.ReadFile("catalog/catalog.yml")
	if err != nil {
		log.Print(err)
		return Catalog{}
	}
	var catalog Catalog
	yaml.Unmarshal(content, &catalog)

	log.Print(catalog)

	return catalog
}


// List return Services provided by the catalog
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

func getProject(service string) (project.APIProject, error) {
	project, err := docker.NewProject(&ctx.Context{
		Context: project.Context{
			ComposeFiles: []string{"catalog/" + service + "/docker-compose.yml"},
			ProjectName:  service,
		},
	}, nil)

	return project, err
}

func Info(writer http.ResponseWriter, request *http.Request) {
	service := mux.Vars(request)["service"]
	project, err := getProject(service)
	if err != nil {
		log.Fatal(err)
	}
	info, err := project.Ps(context.Background())
	log.Print(info)
}
