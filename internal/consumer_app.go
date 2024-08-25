package internal

import (
	"context"
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

	// SubscribeToMessages is a method that will inifinite listens to messages,
	// and send them to callback function.
	go func() {
		if err := c.realtimeService.SubscribeToMessages(context.Background(), func(message string) error {
			c.l.Info("get", slog.String("payload:", message))
			return nil
		}); err != nil {
			c.l.Error("payload doesn't recieved", "err", err)
		}

	}()

	c.l.Info("consumer app started")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	c.l.Info("\nclosing consumer app")

	return nil
}

// Start is creating streams and consumers, on msg it will trigger the
// nats handler, which is start to publish messages in realtime service broadcast channel,
// and SubscribeToMessages will get the broadcasted messages.
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
				c.l.Error("can't publish the message", "err", err)
			}
		}

		if err := msg.Ack(); err != nil {
			c.l.Error("can't acknowledge the message", "err", err)
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
