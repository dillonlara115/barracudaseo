package stream

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dillonlara115/barracudaseo/internal/ai/providers"
	"go.uber.org/zap"
)

// SSEEvent represents a Server-Sent Event
type SSEEvent struct {
	Type string `json:"type"`
	Data string `json:"data,omitempty"`
}

// Handler manages SSE streaming from Gemini to the frontend
type Handler struct {
	logger *zap.Logger
}

// NewHandler creates a new streaming handler
func NewHandler(logger *zap.Logger) *Handler {
	return &Handler{logger: logger}
}

// StreamResponse sets up SSE headers and streams a Gemini response to the client.
// The provider must be a *GeminiProvider with StreamCompletion support.
func (h *Handler) StreamResponse(
	ctx context.Context,
	w http.ResponseWriter,
	provider *providers.GeminiProvider,
	messages []providers.Message,
) (string, error) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return "", fmt.Errorf("streaming not supported")
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	var fullText strings.Builder

	err := provider.StreamCompletion(ctx, messages, func(text string) error {
		fullText.WriteString(text)

		evt := SSEEvent{Type: "chunk", Data: text}
		data, _ := json.Marshal(evt)

		_, writeErr := fmt.Fprintf(w, "data: %s\n\n", data)
		if writeErr != nil {
			return writeErr
		}
		flusher.Flush()
		return nil
	})

	if err != nil {
		errEvt := SSEEvent{Type: "error", Data: err.Error()}
		data, _ := json.Marshal(errEvt)
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
		return fullText.String(), err
	}

	doneEvt := SSEEvent{Type: "done"}
	data, _ := json.Marshal(doneEvt)
	fmt.Fprintf(w, "data: %s\n\n", data)
	flusher.Flush()

	return fullText.String(), nil
}
