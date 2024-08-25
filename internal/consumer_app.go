package internal

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go/jetstream"
	natsHandler "github.com/novychok/goldensbtech/internal/handler/nats"
	"github.com/novychok/goldensbtech/internal/service"
)

type ConsumerApp struct {
	l               *slog.Logger
	natsClient      jetstream.JetStream
	natsHandler     *natsHandler.Handler
	realtimeService service.Realtime
}

func (c *ConsumerApp) ServerConsumerApp() error {

	if err := c.Start(); err != nil {
		c.l.Error("start nats stream failed", "err", err)
	}

	go func() {
		if err := c.realtimeService.SubscribeToMessages(context.Background(), func(message string) error {
			c.l.Info("get", slog.String("payload:", message))
			return nil
		}); err != nil {
			c.l.Error("payload doesn't recieved", "err", err)
		}

	}()

	c.l.Info("consumer app started")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	c.l.Info("\nclosing consumer app")

	return nil
}

func (c *ConsumerApp) Start() error {

	ctx := context.Background()
	stream, err := c.natsClient.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:      "payload",
		Subjects:  []string{"payload.*"},
		Retention: jetstream.InterestPolicy,
	})
	if err != nil {
		return err
	}

	consumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:   "payload_consumer",
		AckPolicy: jetstream.AckExplicitPolicy,
	})
	if err != nil {
		return err
	}

	_, err = consumer.Consume(func(msg jetstream.Msg) {

		switch msg.Subject() {
		case "payload.*":
			if err := c.natsHandler.HandleMessagePub(ctx, msg.Data()); err != nil {
				log.Printf("error to publish the message: %s", err)
			}
		}

		if err := msg.Ack(); err != nil {
			log.Printf("error to acknowledge the message: %s", err)
		}

	})
	if err != nil {
		return err
	}

	return nil
}

func NewConsumerApp(l *slog.Logger,
	natsClient jetstream.JetStream,
	natsHandler *natsHandler.Handler,
	realtimeService service.Realtime) *ConsumerApp {
	return &ConsumerApp{
		l:               l,
		natsClient:      natsClient,
		natsHandler:     natsHandler,
		realtimeService: realtimeService,
	}
}
