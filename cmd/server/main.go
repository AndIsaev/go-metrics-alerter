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
	if err := gen.InitVersion(buildVersion, buildDate, buildCommit); err != nil {
		log.Fatal(err)
	}
	app := New()
	if err := app.StartApp(); err != nil {
		log.Fatal(err)
	}
}
