package main

import (
	"flag"
	"fmt"
	"log"
	"maestro/catalog"
	"net/http"
	"os"
	"path"

	"github.com/jasonlvhit/gocron"
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

	CheckComposeUpdates()

	gocron.Every(1).Minute().Do(UpdateServices)

	<-gocron.Start()

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8888", router))
}
