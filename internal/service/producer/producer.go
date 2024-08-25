package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"math/rand"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/novychok/goldensbtech/internal/entity"
	"github.com/novychok/goldensbtech/internal/service"
)

type srv struct {
	l               *slog.Logger
	natsClient      jetstream.JetStream
	realtimeService service.Realtime
}

// Produce is infinit msg broadcast func, that have 2 cases:
// default for publish a payload, and case to stop the broadcast with stopChan.
func (s *srv) Produce(ctx context.Context, stopChan chan struct{}) error {

	l := s.l.With(slog.String("method", "Produce"))

	go func() {
		for {
			select {
			case <-stopChan:
				s.l.Info("interrupting the payload broadcast")
				return
			default:
				time.Sleep(3 * time.Second)

				payloads, err := s.preparePayload()
				if err != nil {
					l.ErrorContext(ctx, "failed to prepare payloads", "err", err)
					return
				}

				wg := &sync.WaitGroup{}
				for _, payload := range payloads {
					wg.Add(1)
					go func() {
						payloadBytes, err := json.Marshal(payload)
						if err != nil {
							l.ErrorContext(ctx, "failed to marshal the payload", "err", err)
							return
						}

						_, err = s.natsClient.Publish(context.Background(), "payload.*", payloadBytes)
						if err != nil {
							l.ErrorContext(ctx, "failed to publish the payload", "err", err)
							return
						}
						wg.Done()
					}()
				}
				wg.Wait()
			}
		}
	}()

	return nil
}

// preparePayload is preparing concurrent Payloads slice.
func (s *srv) preparePayload() ([]entity.Payload, error) {

	numPayloads := rand.Intn(6) + 5

	payloads := make([]entity.Payload, 0, numPayloads)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 0; i < numPayloads; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			payload := entity.Payload{
				ID:    fmt.Sprintf("payload%d", i+1),
				Count: rand.Intn(100) + 1,
			}

			mu.Lock()
			payloads = append(payloads, payload)
			mu.Unlock()
		}(i)
	}

	wg.Wait()

	return payloads, nil
}

func New(l *slog.Logger,
	natsClient jetstream.JetStream,
	realtimeService service.Realtime) service.Producer {
	return &srv{
		l:               l,
		natsClient:      natsClient,
		realtimeService: realtimeService,
	}
}
