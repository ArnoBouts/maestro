package main

import (
	"io/ioutil"
	"log"
	"os"
	"net/http"

        "gopkg.in/yaml.v2"

        "golang.org/x/net/context"

        "github.com/docker/libcompose/docker"
        "github.com/docker/libcompose/docker/ctx"
        "github.com/docker/libcompose/project"
        "github.com/docker/libcompose/project/options"

        "github.com/gorilla/mux"
)

type maestro struct {
	Services map[string](*Service)
}

type Service struct {
	Name string
	Enable bool
}

var m = new(maestro)

func Load() {

	// load from file
        content, err := ioutil.ReadFile(workdir + "/services/services.yml")
	if err != nil {
                log.Println("Unable to read services file")
		return
        }

	yaml.Unmarshal(content, &m)
	log.Print("M : ")
	log.Println(m)

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

func Save() {

	log.Print("M : ")
	log.Println(m)

	content, _ := yaml.Marshal(&m)
	ioutil.WriteFile(workdir + "/services/services.yml", content, 0644)
}

func add(name string) {
	log.Println("Install service '" + name + "'")

        compose, err := ioutil.ReadFile(workdir + "/catalog/" + name + "/docker-compose.yml")
        if err != nil {
                log.Fatal(err)
        }

	//app := c.Apps[name]

	// copy the docker-compose.yml
	err = os.Mkdir(workdir + "/services/" + name, 0777)
        if err != nil {
                log.Fatal(err)
        }
	err = ioutil.WriteFile(workdir + "/services/" + name + "/docker-compose.yml", []byte(compose), 0644)
        if err != nil {
                log.Fatal(err)
        }

	// save parameters

	// add service to maestro
	var service Service
	service.Name = name
	if(m.Services == nil) {
		m.Services = make(map[string](*Service))
	}
	m.Services[name] = &service
	Save()
	
	// up compose
	service.Up()
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

func (service *Service) Info() {

        project, err := getProject(service.Name)

        if err != nil {
                log.Fatal(err)
        }

        info, err := project.Ps(context.Background())

        if err != nil {
                log.Println(err)
        }
	log.Print("Infos : ")
	log.Println(info)
}

//StartService call start method on compose service
func (service *Service) Start() {

        project, err := getProject(service.Name)

        if err != nil {
                log.Fatal(err)
        }

        err = project.Start(context.Background())

        if err != nil {
                log.Println(err)
        }
}

//StopService call stop method on compose service
func (service *Service) Stop() {

        project, err := getProject(service.Name)

        if err != nil {
                log.Fatal(err)
        }

        err = project.Stop(context.Background(), 10000)

        if err != nil {
                log.Fatal(err)
        }
}

//UpService call up method on compose service
func (service *Service) Up() {

        project, err := getProject(service.Name)

        if err != nil {
                log.Fatal(err)
        }

        err = project.Up(context.Background(), options.Up{})

        if err != nil {
                log.Fatal(err)
        }
}

//DownService call down method on compose service
func (service *Service) Down() {

        project, err := getProject(service.Name)

        if err != nil {
                log.Fatal(err)
        }

        err = project.Down(context.Background(), options.Down{})

        if err != nil {
                log.Fatal(err)
        }
}

//InfoService call info method on compose service
func InfoService(writer http.ResponseWriter, request *http.Request) {

	log.Println(m)

	service := m.Services[mux.Vars(request)["service"]]
	service.Info()
}

//StartService call start method on compose service
func StartService(writer http.ResponseWriter, request *http.Request) {

	service := m.Services[mux.Vars(request)["service"]]
	service.Start()
	service.Enable = true
	log.Println(service)
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
	add(name)
}
