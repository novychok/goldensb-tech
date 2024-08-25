package goldensbtech

import "github.com/novychok/goldensbtech/internal"

type Goldensb interface {
	ServeProducerApp() error
	ServeConsumerApp() error
}

type goldensb struct{}

func (g *goldensb) ServeProducerApp() error {

	app, cleanup, err := internal.InitProducerApp()
	if err != nil {
		return err
	}
	defer cleanup()

	err = app.ServerProducerApp()
	if err != nil {
		return err
	}

	return nil
}

func (g *goldensb) ServeConsumerApp() error {

	app, cleanup, err := internal.InitConsumerApp()
	if err != nil {
		return err
	}
	defer cleanup()

	err = app.ServerConsumerApp()
	if err != nil {
		return err
	}

	return nil
}

func New() Goldensb {
	return &goldensb{}
}
