package main

import (
	"log"
	"net/http"
	"maestro/catalog"
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

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8888", router))
}
