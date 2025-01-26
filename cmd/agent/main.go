package main

import (
	"log"

	"github.com/AndIsaev/go-metrics-alerter/internal/service/gen"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	log.Println("start app")
	if err := gen.InitVersion(buildVersion, buildDate, buildCommit); err != nil {
		log.Fatal(err)
	}
	app := New()
	app.StartApp()
}
