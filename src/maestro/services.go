package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v2"

	"golang.org/x/net/context"

	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/docker/ctx"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/project/options"

	"github.com/gorilla/mux"

	"maestro/catalog"
)

type maestro struct {
	Services map[string](*Service)
}

// Service is an installed service
type Service struct {
	Name     string
	Enable   bool
	Checksum string
	Params   map[string](string)
}

var m = new(maestro)

// Load installed service descriptor file
func Load() {

	log.Println("Loading services descriptor file")

	// load from file
	content, err := ioutil.ReadFile(workdir + "/services/services.yml")
	if err != nil {
		log.Println("Unable to read services descriptor file")
		return
	}

	yaml.Unmarshal(content, &m)

	if m.Services == nil {
		m.Services = make(map[string](*Service))
	}

	// start all that should
	for name, service := range m.Services {
		service.Name = name
		if service.Enable {
			service.start()
		}
		log.Println(service)
	}
}

// CheckComposeUpdates check for each service is the compose file was updated
func CheckComposeUpdates() {

	// for all enabled services, check with sha256 if compose was updated in the catalog
	for name, service := range m.Services {
		if service.Enable {
			sha, _ := catalog.ComposeSha256(name)
			if service.Checksum != sha {
				log.Println(name + " compose file need to be updated")
			}
		}
	}
}

// Save save the services descriptor file
func Save() {

	content, _ := yaml.Marshal(&m)
	ioutil.WriteFile(workdir+"/services/services.yml", content, 0644)
}

func add(name string, params map[string](string)) error {
	log.Println("Install service '" + name + "'")

	// create service
	var service Service
	service.Name = name
	service.Params = params

	compose, err := catalog.ComposeFile(service.Name)
	if err != nil {
		return err
	}

	// write service compose file
	err = service.writeCompose(compose)
	if err != nil {
		return err
	}

	// add service to maestro
	if m.Services == nil {
		m.Services = make(map[string](*Service))
	}
	m.Services[name] = &service
	Save()

	// up compose
	return service.up()
}

func (service *Service) writeCompose(compose string) error {

	for k, v := range service.Params {
		compose = strings.Replace(compose, "{{"+k+"}}", v, -1)
	}

	if err := os.Mkdir(workdir+"/services/"+service.Name, 0777); err != nil {
		return err
	}
	return ioutil.WriteFile(workdir+"/services/"+service.Name+"/docker-compose.yml", []byte(compose), 0644)
}

func getProject(service string) (project.APIProject, error) {
	project, err := docker.NewProject(&ctx.Context{
		Context: project.Context{
			ComposeFiles: []string{workdir + "/services/" + service + "/docker-compose.yml"},
			ProjectName:  service,
		},
	}, nil)
	return project, err
}

func (service *Service) info() (project.InfoSet, error) {

	project, err := getProject(service.Name)
	if err != nil {
		return nil, err
	}

	return project.Ps(context.Background())
}

func (service *Service) start() error {

	log.Printf("Start service '%s'\n", service.Name)

	project, err := getProject(service.Name)
	if err != nil {
		return err
	}

	err = project.Start(context.Background())

	if err != nil {
		log.Printf("Service '%s' starting failed\n", service.Name)
	} else {
		log.Printf("Service '%s' started\n", service.Name)
	}

	return err
}

func (service *Service) stop() error {

	project, err := getProject(service.Name)
	if err != nil {
		return err
	}

	return project.Stop(context.Background(), 10000)
}

func (service *Service) up() error {

	project, err := getProject(service.Name)
	if err != nil {
		return err
	}

	return project.Up(context.Background(), options.Up{})
}

func (service *Service) down() error {

	project, err := getProject(service.Name)
	if err != nil {
		return err
	}

	return project.Down(context.Background(), options.Down{})
}

func (service *Service) updateCompose() error {

	if err := service.down(); err != nil {
		return err
	}

	compose, err := catalog.ComposeFile(service.Name)
	if err != nil {
		return err
	}

	// write service compose file
	if err := service.writeCompose(compose); err != nil {
		return err
	}

	return service.up()
}

// InfoService Resource that return provided service infos
func InfoService(writer http.ResponseWriter, request *http.Request) {

	service := m.Services[mux.Vars(request)["service"]]
	info, err := service.info()
	if err != nil {
		http.Error(writer, err.Error(), 500)
		return
	}

	payload, err := json.Marshal(info)
	if err != nil {
		http.Error(writer, err.Error(), 500)
		return
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.Write(payload)
}

// StartService Resource that start the provided service
func StartService(writer http.ResponseWriter, request *http.Request) {

	service := m.Services[mux.Vars(request)["service"]]
	service.start()
	service.Enable = true
	Save()
}

// StopService Resource that stop the provided service
func StopService(writer http.ResponseWriter, request *http.Request) {

	service := m.Services[mux.Vars(request)["service"]]
	service.stop()
	service.Enable = false
	Save()
}

// UpService Resource that up the provided service
func UpService(writer http.ResponseWriter, request *http.Request) {

	var service Service
	service.Name = mux.Vars(request)["service"]

	service.up()
}

// DownService Resource that down the provided service
func DownService(writer http.ResponseWriter, request *http.Request) {

	var service Service
	service.Name = mux.Vars(request)["service"]

	service.down()
}

// AddService Resource that install the provided service
func AddService(writer http.ResponseWriter, request *http.Request) {
	name := mux.Vars(request)["service"]

	var params map[string](string)
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&params); err != nil {
		http.Error(writer, err.Error(), 500)
		return
	}
	defer request.Body.Close()

	if err := add(name, params); err != nil {
		http.Error(writer, err.Error(), 500)
		return
	}
}
