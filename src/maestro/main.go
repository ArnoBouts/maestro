package main

import (
	"flag"
	"fmt"
	"log"
	"maestro/catalog"
	"net/http"
	"os"
	"path"
)

var workdir string

func main() {

	fmt.Println(os.Args)

	restart := flag.Bool("restart", false, "maestro app must be restarted")

	flag.Parse()

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	workdir = path.Dir(ex)

	catalog.Load(workdir)

	if *restart {
		Restart()
		return
	}

	Load()

	InstallRequired()

	CheckImageToUpdate()

	CheckComposeUpdates()

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8888", router))
}
