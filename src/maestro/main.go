package main

import (
	"log"
	"net/http"
	"maestro/catalog"
)

var c catalog.Catalog

func main() {

	c = catalog.Load()

	Load()

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8888", router))
}
