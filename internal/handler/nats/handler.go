package nats

import (
	"context"

	"github.com/novychok/goldensbtech/internal/service"
)

type Handler struct {
	realtimeService service.Realtime
}

func (h *Handler) HandleMessagePub(ctx context.Context, data []byte) error {
	return h.realtimeService.PublishMessage(ctx, string(data))
}

func New(realtimeService service.Realtime) *Handler {
	return &Handler{
		realtimeService: realtimeService,
	}
}
