// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package internal

import (
	nats2 "github.com/novychok/goldensbtech/internal/handler/nats"
	producer2 "github.com/novychok/goldensbtech/internal/handler/producer"
	"github.com/novychok/goldensbtech/internal/pkg/log"
	"github.com/novychok/goldensbtech/internal/pkg/nats"
	"github.com/novychok/goldensbtech/internal/service/producer"
	"github.com/novychok/goldensbtech/internal/service/realtime"
)

// Injectors from wire.go:

func InitProducerApp() (*ProducerApp, func(), error) {
	logger := log.New()
	jetStream, cleanup, err := nats.New()
	if err != nil {
		return nil, nil, err
	}
	serviceRealtime := realtime.New()
	serviceProducer := producer.New(logger, jetStream, serviceRealtime)
	handler := producer2.New(serviceProducer, serviceRealtime)
	producerApp := NewProducerApp(logger, jetStream, handler)
	return producerApp, func() {
		cleanup()
	}, nil
}

func InitConsumerApp() (*ConsumerApp, func(), error) {
	logger := log.New()
	jetStream, cleanup, err := nats.New()
	if err != nil {
		return nil, nil, err
	}
	serviceRealtime := realtime.New()
	handler := nats2.New(serviceRealtime)
	consumerApp := NewConsumerApp(logger, jetStream, handler, serviceRealtime)
	return consumerApp, func() {
		cleanup()
	}, nil
}
