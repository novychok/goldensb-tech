package producer

import (
	"encoding/json"
	"net/http"

	"github.com/novychok/goldensbtech/internal/service"
)

type Handler struct {
	producerService service.Producer
}

func (h *Handler) SendMessages(w http.ResponseWriter, r *http.Request) {

	if err := h.producerService.Produce(r.Context()); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"status": "Internal Server Error"})
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]string{
		"status": "the payload was sended well"})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return
	}
}

func New(producerService service.Producer) *Handler {
	return &Handler{
		producerService: producerService,
	}
}
