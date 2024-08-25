//go:build wireinject
// +build wireinject

package internal

import (
	"github.com/google/wire"
	"github.com/novychok/goldensbtech/internal/handler/producer"

	natsHandler "github.com/novychok/goldensbtech/internal/handler/nats"
	"github.com/novychok/goldensbtech/internal/pkg"
	"github.com/novychok/goldensbtech/internal/service/service"
)

func InitProducerApp() (*ProducerApp, func(), error) {
	wire.Build(
		pkg.Set,
		service.Set,
		producer.New,
		NewProducerApp,
	)
	return nil, nil, nil
}

func InitConsumerApp() (*ConsumerApp, func(), error) {
	wire.Build(
		pkg.Set,
		service.Set,
		natsHandler.New,
		NewConsumerApp,
	)
	return nil, nil, nil
}
