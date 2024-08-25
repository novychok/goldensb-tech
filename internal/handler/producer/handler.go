package producer

import (
	"encoding/json"
	"net/http"

	"github.com/novychok/goldensbtech/internal/service"
)

type Handler struct {
	producerService service.Producer
	realtimeService service.Realtime
	stopChan        chan struct{}
}

func (h *Handler) SendMessages(w http.ResponseWriter, r *http.Request) {

	if err := h.producerService.Produce(r.Context(), h.stopChan); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"InternalServerError": "error to broadcast the payload"})
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]string{
		"status": "the payload was sended well"})
}

func (h *Handler) StopSendingMessages(w http.ResponseWriter, r *http.Request) {

	h.stopChan <- struct{}{}

	writeJSON(w, http.StatusAccepted, map[string]string{
		"status": "stop sending payload"})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return
	}
}

func New(producerService service.Producer,
	realtimeService service.Realtime) *Handler {
	return &Handler{
		producerService: producerService,
		realtimeService: realtimeService,
		stopChan:        make(chan struct{}),
	}
}
