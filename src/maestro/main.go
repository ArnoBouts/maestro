package main

import (
	"log"
	"maestro/catalog"
	"net/http"
	"os"
	"path"
)

var workdir string

func main() {

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	workdir = path.Dir(ex)

	catalog.Load(workdir)

	Load()

	InstallRequired()

	CheckImageToUpdate()

	CheckComposeUpdates()

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8888", router))
}
