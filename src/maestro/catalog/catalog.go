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
	"github.com/docker/libcompose/project/options"

	"github.com/gorilla/mux"
)

// Catalog define the catalog
type Catalog struct {
	Services map[string]Service
}

// Service define a service provided by the catalogs
type Service struct {
	DisplayName string `yaml:"display_name"`
	Required    bool   `yaml:"required"`
	//Params []Param
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

//StartService call start method on compose service
func StartService(writer http.ResponseWriter, request *http.Request) {

	service := mux.Vars(request)["service"]

	project, err := getProject(service)

	if err != nil {
		log.Fatal(err)
	}

	err = project.Start(context.Background())

	if err != nil {
		log.Fatal(err)
	}
}

//StopService call stop method on compose service
func StopService(writer http.ResponseWriter, request *http.Request) {

	service := mux.Vars(request)["service"]

	project, err := getProject(service)

	if err != nil {
		log.Fatal(err)
	}

	err = project.Stop(context.Background(), 10000)

	if err != nil {
		log.Fatal(err)
	}
}

//UpService call up method on compose service
func UpService(writer http.ResponseWriter, request *http.Request) {

	service := mux.Vars(request)["service"]

	project, err := getProject(service)

	if err != nil {
		log.Fatal(err)
	}

	err = project.Up(context.Background(), options.Up{})

	if err != nil {
		log.Fatal(err)
	}
}

//DownService call down method on compose service
func DownService(writer http.ResponseWriter, request *http.Request) {

	service := mux.Vars(request)["service"]

	project, err := getProject(service)

	if err != nil {
		log.Fatal(err)
	}

	err = project.Down(context.Background(), options.Down{})

	if err != nil {
		log.Fatal(err)
	}
}
