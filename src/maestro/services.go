package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"

	"golang.org/x/net/context"

	"github.com/docker/docker/client"
	composeclient "github.com/docker/libcompose/docker/client"
	"github.com/docker/libcompose/docker/container"
	"github.com/docker/libcompose/docker/ctx"
	"github.com/docker/libcompose/docker/image"
	"github.com/docker/libcompose/labels"
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
	}
}

// InstallRequired install required services
func InstallRequired() {

	log.Println("Install required services")

	for _, name := range catalog.GetRequiredApps() {

		if _, contains := m.Services[name]; !contains {
			err := add(name, make(map[string](string)))
			if err != nil {
				log.Println(err.Error())
			}
		}
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
				err := service.down()
				if err != nil {
					log.Printf("Enable to down service %s : %s", service.Name, err.Error())
					continue
				}

				//override compose file
				p, err := service.computeParams(service.Params)
				if err != nil {
					log.Printf("Enable to compute params for the service %s : %s", service.Name, err.Error())
					continue
				}
				service.Params = p

				err = service.configure()
				if err != nil {
					log.Printf("Enable to configure service %s : %s", service.Name, err.Error())
					continue
				}

				err = service.up()
				if err != nil {
					log.Printf("Enable to up service %s : %s", service.Name, err.Error())
					continue
				}
			
				log.Println(name + " compose file updated")
			}
		}
	}
}

// PullServices pull all services images
func PullServices() {

	for _, service := range m.Services {
		if service.Enable {
			service.pull()
		}
	}
}

func CheckImageToUpdate() {

	for _, service := range m.Services {
		if service.Enable {

			if err := service.checkImageToUpdate(); err != nil {
				log.Printf("Unable to uptade service %s : %s", service.Name, err.Error())
			}
		}
	}
}

func (service *Service) checkImageToUpdate() error {

	uptodate := true

	p, err := getProject(service.Name)
	if err != nil {
		return err
	}

	for _, name := range p.ServiceConfigs.Keys() {
		containers, _ := collectContainers(context.Background(), p, name)

		for _, c := range containers {
			outOfSync, err := outOfSync(context.Background(), c, p, name)
			if err != nil {
				return err
			}

			if outOfSync {
				log.Printf("%s is out of sync", name)
				uptodate = false
			}
		}
	}

	if !uptodate {
		if err := service.update(); err != nil {
			return err
		}
	}

	return nil
}

func collectContainers(ctx context.Context, p project.Project, service string) ([]*container.Container, error) {
	client, _ := composeclient.Create(composeclient.Options{})
	containers, err := container.ListByFilter(ctx, client, labels.SERVICE.Eq(service), labels.PROJECT.Eq(p.Name))
	if err != nil {
		return nil, err
	}

	result := []*container.Container{}

	for _, cont := range containers {
		c, err := container.New(ctx, client, cont.ID)
		if err != nil {
			return nil, err
		}
		result = append(result, c)
	}

	return result, nil
}

func outOfSync(ctx context.Context, c *container.Container, p project.Project, service string) (bool, error) {

	conf, ok := p.GetServiceConfig(service)
	if !ok {
		return false, fmt.Errorf("Failed to find service: %s", service)
	}

	if c.ImageConfig() != conf.Image {
		log.Printf("Images for %s do not match %s!=%s", c.Name(), c.ImageConfig(), conf.Image)
		return true, nil
	}

	// TODO : issue when trying to check service hash. conf pb ?
	//
	//expectedHash := config.GetServiceHash(service, conf)
	//if c.Hash() != expectedHash {
	//	log.Printf("Hashes for %s do not match %s!=%s", c.Name(), c.Hash(), expectedHash)
	//	return true, nil
	//}

	cli, _ := composeclient.Create(composeclient.Options{})

	image, err := image.InspectImage(ctx, cli, c.ImageConfig())
	if err != nil {
		if client.IsErrImageNotFound(err) {
			log.Printf("Image %s do not exist, do not know if it's out of sync", c.Image())
			return false, nil
		}
		return false, err
	}

	return image.ID != c.Image(), err
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
	p, err := service.computeParams(params)
	if err != nil {
		return err
	}
	service.Params = p

	err = service.configure()
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
	err = service.up()
	if err != nil {
		return err
	}
	service.Enable = true
	Save()
	return nil
}

func (service *Service) configure() error {

	sha, err := catalog.ComposeSha256(service.Name)
	service.Checksum = sha
	if err != nil {
		return err
	}

	compose, err := catalog.ComposeFile(service.Name)
	if err != nil {
		return err
	}

	// write service compose file
	err = service.writeCompose(compose)
	if err != nil {
		return err
	}

	return nil
}

func (service *Service) computeParams(params map[string](string)) (map[string](string), error) {
	result := make(map[string](string))

	catalogParams := catalog.GetServiceParams(service.Name)

	if(catalogParams != nil) {
		for p, _ := range catalogParams {
			v, err := service.getParamValue(p)
			if err != nil {
				return nil, err
			}
			result[p] = v
		}
	}

	return result, nil
}

func (service *Service) writeCompose(compose string) error {

	r := regexp.MustCompile(`{{([^}]*)}}`)

	params := r.FindAllStringSubmatch(compose, -1)

	for _, param := range params {

		// Search for param value
		val, err := service.getParamValue(param[1])
		if err != nil {
			return err
		}

		compose = strings.Replace(compose, param[0], val, -1)
	}

	if err := os.MkdirAll(workdir+"/services/"+service.Name, 0777); err != nil {
		return err
	}
	return ioutil.WriteFile(workdir+"/services/"+service.Name+"/docker-compose.yml", []byte(compose), 0644)
}

func (service *Service) getParamValue(param string) (string, error) {

	if val, founded := service.Params[param]; founded {
		return val, nil
	}
	if val, founded := os.LookupEnv(param); founded {
		return val, nil
	}
	if val, founded := catalog.GetServiceParam(service.Name, param); founded {
		return val, nil
	}

	// If not value found and param required, return an error
	return "", fmt.Errorf("Undefined required param value for '%s'", param)
}

func getProject(service string) (project.Project, error) {
	p, err := NewProject(&ctx.Context{
		Context: project.Context{
			ComposeFiles: []string{workdir + "/services/" + service + "/docker-compose.yml"},
			ProjectName:  service,
		},
	}, nil)
	return *p, err
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

func (service *Service) pull() error {

	project, err := getProject(service.Name)
	if err != nil {
		return err
	}

	return project.Pull(context.Background())
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

func (service *Service) update() error {
	if err := service.pull(); err != nil {
		return err
	}

	updater := catalog.GetUpdater(service.Name)

	if updater == "" {
		return service.up()

	} else {
		s, founded := m.Services[updater]
		if !founded {
			return add(updater, make(map[string](string)))
		} else {
			if err := s.pull(); err != nil {
				return err
			}
			return s.up()
		}
	}
}

// UpdateService Resource that update the provided service
func UpdateService(writer http.ResponseWriter, request *http.Request) {

	var service Service
	service.Name = mux.Vars(request)["service"]

	if err := service.update(); err != nil {
		http.Error(writer, err.Error(), 500)
	}
}

func Restart() {
	var service Service
	service.Name = "maestro"

	service.up()
}

func UpdateServices() {

	PullServices()

	CheckImageToUpdate()
}
