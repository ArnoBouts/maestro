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

	c.Workdir = workdir
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
