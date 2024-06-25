package main

import (
	"fmt"
	"github.com/AndIsaev/go-metrics-alerter/internal/logger"
	"github.com/AndIsaev/go-metrics-alerter/internal/service"
	"net/http"
)

func run() error {
	if err := logger.Initialize(); err != nil {
		return err
	}

	config := service.NewServerConfig()

	fmt.Println("Running server on", config.Address)
	return http.ListenAndServe(config.Address, config.Route)

}

func main() {

	if err := run(); err != nil {
		println(err)
	}
}
