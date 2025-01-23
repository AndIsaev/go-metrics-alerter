package main

import (
	"context"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/gen"
	"log"
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
	ctx := context.Background()
	err := app.StartApp(ctx)

	defer func() {
		app.Shutdown()
		if err != nil {
			log.Fatalf("close process with error: %s\n", err.Error())
		}
	}()
}
