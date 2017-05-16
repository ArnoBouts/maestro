package catalog

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

// Catalog define the catalog
type Catalog struct {
	Workdir string
	Apps    map[string]App `yaml:"services"`
}

// App define an app provided by the catalogs
type App struct {
	DisplayName string           `yaml:"display_name"`
	Required    bool             `yaml:"required"`
	Params      map[string]Param `yaml:"params"`
}

// Param define a parameter of an app
type Param struct {
	Required bool   `yaml:"required"`
	Default  string `yaml:"default"`
}

var c Catalog

func Load(workdir string) {
	content, err := ioutil.ReadFile(workdir + "/catalog/catalog.yml")
	if err != nil {
		log.Print(err)
		return
	}
	yaml.Unmarshal(content, &c)

	c.Workdir = workdir
	log.Println(c)
}

func GetRequiredApps() []string {

	var requiredApps []string

	for name, app := range c.Apps {
		if app.Required {
			requiredApps = append(requiredApps, name)
		}
	}

	return requiredApps
}

func GetServiceParam(service string, param string) (string, bool) {

	if s, f := c.Apps[service]; f {
		if p, f := s.Params[param]; f {
			if p.Required && p.Default == "" {
				return "", false
			}
			return p.Default, true
		}
	}

	log.Println("Pas trouvé le param")
	return "", false
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

func ComposeFile(name string) (string, error) {
	c, err := ioutil.ReadFile(c.Workdir + "/catalog/" + name + "/docker-compose.yml")
	return string(c), err
}

func ComposeSha256(name string) (string, error) {

	f, err := os.Open(c.Workdir + "/catalog/" + name + "/docker-compose.yml")
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
