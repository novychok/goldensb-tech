package nats

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func New() (jetstream.JetStream, func(), error) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, nil, err
	}

	js, err := jetstream.New(nc)
	if err != nil {
		return nil, nil, err
	}

	cleaner := func() {
		nc.Close()
	}

	return js, cleaner, nil
}
