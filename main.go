package main

import (
	"github.com/deuterium34/go-rave/app"
)

func main() {
	app, err := app.NewApp()
	if err != nil {
		panic(err)
	}
	defer app.Close()

	if err = app.Start(); err != nil {
		panic(err)
	}

	if err := <-app.CloseCh; err != nil {
		panic(err)
	}
}
