package main

import (
	"log"
	"net/http"
	"maestro/catalog"
)

var c catalog.Catalog

func main() {

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := path.Dir(ex)
	log.Print(exPath)

	c = catalog.Load()

	Load()

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8888", router))
}
