package main

import "github.com/novychok/goldensbtech"

func main() {
	app := goldensbtech.New()
	if err := app.ServeProducerApp(); err != nil {
		panic(err)
	}
}
