package producer

import (
	"context"
	"log/slog"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/novychok/goldensbtech/internal/service"
)

type srv struct {
	l          *slog.Logger
	natsClient jetstream.JetStream
}

func (s *srv) Produce(ctx context.Context) error {

	l := s.l.With(slog.String("method", "Produce"))

	// payload :=

	_, err := s.natsClient.Publish(context.Background(), "payload.*", []byte("miss monyque "))
	if err != nil {
		l.ErrorContext(ctx, "failed to publish the payload", "err", err)
		return err
	}

	return nil
}

func New(l *slog.Logger,
	natsClient jetstream.JetStream) service.Producer {
	return &srv{
		l:          l,
		natsClient: natsClient,
	}
}
