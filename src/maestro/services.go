package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"net/http"
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

type Service struct {
	Name string
	Enable bool
	Checksum string
	Params map[string](string)
}

var m = new(maestro)

func Load() {

        log.Println("Loading services descriptor file")

	// load from file
        content, err := ioutil.ReadFile(workdir + "/services/services.yml")
	if err != nil {
                log.Println("Unable to read services descriptor file")
		return
        }

	yaml.Unmarshal(content, &m)

	if(m.Services == nil) {
		m.Services = make(map[string](*Service))
	}

	// start all that should
	for name, service := range m.Services {
		service.Name = name
		if(service.Enable) {
			service.Start()
		}
		log.Println(service)
	}
}

func CheckComposeUpdates() {
	
	// for all enabled services, check with sha256 if compose was updated in the catalog
	for name, service := range m.Services {
		if(service.Enable) {
			sha, _ := catalog.ComposeSha256(name)
			if service.Checksum != sha {
				log.Println(name + " compose file need to be updated")
			}
		}
	}
}

func Save() {

	content, _ := yaml.Marshal(&m)
	ioutil.WriteFile(workdir + "/services/services.yml", content, 0644)
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
	if(m.Services == nil) {
		m.Services = make(map[string](*Service))
	}
	m.Services[name] = &service
	Save()
	
	// up compose
	return service.Up()
}

func (service *Service) writeCompose(compose string) error {

	for k, v := range service.Params {
		compose = strings.Replace(compose, "{{" + k + "}}", v, -1)
	}

	if err := os.Mkdir(workdir + "/services/" + service.Name, 0777); err != nil {
		return err
	}
	return ioutil.WriteFile(workdir + "/services/" + service.Name + "/docker-compose.yml", []byte(compose), 0644)
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

func (service *Service) Info() (project.InfoSet, error) {

        project, err := getProject(service.Name)
	if err != nil {
		return nil, err
        }

        return project.Ps(context.Background())
}

//StartService call start method on compose service
func (service *Service) Start() error {

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

//StopService call stop method on compose service
func (service *Service) Stop() error {

        project, err := getProject(service.Name)
	if err != nil {
		return err
        }

        return project.Stop(context.Background(), 10000)
}

//UpService call up method on compose service
func (service *Service) Up() error {

        project, err := getProject(service.Name)
	if err != nil {
		return err
	}

        return project.Up(context.Background(), options.Up{})
}

//DownService call down method on compose service
func (service *Service) Down() error {

	project, err := getProject(service.Name)
	if err != nil {
		return err
	}

        return project.Down(context.Background(), options.Down{})
}

func (service *Service) UpdateCompose() error {

	if err := service.Down(); err != nil {
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

	return service.Up()
}

//InfoService call info method on compose service
func InfoService(writer http.ResponseWriter, request *http.Request) {

	service := m.Services[mux.Vars(request)["service"]]
	info, err := service.Info()
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

//StartService call start method on compose service
func StartService(writer http.ResponseWriter, request *http.Request) {

	service := m.Services[mux.Vars(request)["service"]]
	service.Start()
	service.Enable = true
	Save()
}

//StopService call stop method on compose service
func StopService(writer http.ResponseWriter, request *http.Request) {

	service := m.Services[mux.Vars(request)["service"]]
	service.Stop()
	service.Enable = false
	Save()
}

//UpService call up method on compose service
func UpService(writer http.ResponseWriter, request *http.Request) {

	var service Service
        service.Name = mux.Vars(request)["service"]

	service.Up()
}

//DownService call down method on compose service
func DownService(writer http.ResponseWriter, request *http.Request) {

	var service Service
        service.Name = mux.Vars(request)["service"]

	service.Down()
}

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
