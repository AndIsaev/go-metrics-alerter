package main

import "log"

func main() {
	log.Println("start app")
	app := New()
	app.StartApp()
}
