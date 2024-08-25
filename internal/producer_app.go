package internal

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/novychok/goldensbtech/internal/handler/producer"
)

type ProducerApp struct {
	l               *slog.Logger
	natsClient      jetstream.JetStream
	producerHandler *producer.Handler
}

func (p *ProducerApp) ServerProducerApp() error {

	port := flag.Int("port", 4444, "port for producer app")
	flag.Parse()

	http.HandleFunc("GET /produce", p.producerHandler.SendMessages)
	http.HandleFunc("GET /produce-stop", p.producerHandler.StopSendingMessages)

	p.l.Info("producer app listening",
		slog.Int("port", *port))

	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		return err
	}

	return nil
}

func NewProducerApp(l *slog.Logger,
	natsClient jetstream.JetStream,
	producerHandler *producer.Handler) *ProducerApp {
	return &ProducerApp{
		l:               l,
		natsClient:      natsClient,
		producerHandler: producerHandler,
	}
}
