package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/AndIsaev/go-metrics-alerter/internal/logger"
	"github.com/AndIsaev/go-metrics-alerter/internal/service"
)

func run() error {
	if err := logger.Initialize(); err != nil {
		return err
	}
	config := service.NewServerConfig()
	app := NewServerApp(config)

	if app.DbConn != nil {
		defer app.DbConn.Close(context.Background())
	}

	fmt.Println("Running server on", config.Address)
	return http.ListenAndServe(config.Address, app.Route)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
