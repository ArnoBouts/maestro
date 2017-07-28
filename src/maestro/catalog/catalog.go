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
	Updater     string           `yaml:"updater"`
	Required    bool             `yaml:"required"`
	Params      map[string]Param `yaml:"params"`
	LdapGroup   string           `yaml:"ldap_group"`
	Install     []Command        `yaml:"install"`
}

// Param define a parameter of an app
type Param struct {
	Required bool   `yaml:"required"`
	Default  string `yaml:"default"`
}

// Command define a docker run command
type Command struct {
	Service string   `yaml:"service"`
	Command []string `yaml:"command"`
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

func GetApp(name string) *App {

	if a, f := c.Apps[name]; f {
		return &a
	}

	return nil
}

func GetServiceParams(service string) map[string]Param {

	if s, f := c.Apps[service]; f {
		return s.Params
	}

	return nil
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

	return "", false
}

func GetUpdater(service string) string {
	return c.Apps[service].Updater
}

func GetLdapGroup(service string) string {
	return c.Apps[service].LdapGroup
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
