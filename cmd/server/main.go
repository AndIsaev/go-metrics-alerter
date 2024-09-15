package main

import (
	"context"
	"log"
)

func main() {
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
