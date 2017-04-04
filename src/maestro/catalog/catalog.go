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

func getProject(service string) (project.APIProject, error) {
	project, err := docker.NewProject(&ctx.Context{
		Context: project.Context{
			ComposeFiles: []string{"/catalog/" + service + "/docker-compose.yml"},
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
