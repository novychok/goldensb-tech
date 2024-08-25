package main

import "github.com/novychok/goldensbtech"

func main() {
	app := goldensbtech.New()
	if err := app.ServeConsumerApp(); err != nil {
		panic(err)
	}
}
