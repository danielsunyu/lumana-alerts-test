package main

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type EventType string

const (
	EventPlug   EventType = "plug"
	EventUnplug EventType = "unplug"
)

// maxBodyBytes caps the accepted request body size.
const maxBodyBytes = 1 << 20 // 1 MiB

// Handler receives Lumana HTTP alerts.
type Handler struct {
	log *slog.Logger
}

// NewRouter wires up the chi router.
func NewRouter(log *slog.Logger) http.Handler {
	h := &Handler{log: log}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/healthz"))

	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("Hello, World!"))
	})

	r.Route("/alerts", func(r chi.Router) {
		r.Post("/plug", h.handlePlug)
		r.Post("/unplug", h.handleUnplug)
	})

	return r
}

// handlePlug handles the plug event.
func (h *Handler) handlePlug(w http.ResponseWriter, r *http.Request) {
	payload, ok := h.receive(w, r, EventPlug)
	if !ok {
		return
	}
	_ = payload // TODO: plug-specific processing.
	writeOK(w)
}

// handleUnplug handles the unplug event.
func (h *Handler) handleUnplug(w http.ResponseWriter, r *http.Request) {
	payload, ok := h.receive(w, r, EventUnplug)
	if !ok {
		return
	}
	_ = payload // TODO: unplug-specific processing.
	writeOK(w)
}

// receive reads and logs the alert payload. On failure it writes the error
// response and returns false.
func (h *Handler) receive(w http.ResponseWriter, r *http.Request, event EventType) (any, bool) {
	payload, err := readPayload(w, r)
	if err != nil {
		if errors.As(err, new(*http.MaxBytesError)) {
			http.Error(w, "body too large", http.StatusRequestEntityTooLarge)
		} else {
			http.Error(w, "failed to read body", http.StatusBadRequest)
		}
		return nil, false
	}

	h.log.Info("lumana alert received",
		"event", event,
		"request_id", middleware.GetReqID(r.Context()),
		"remote", r.RemoteAddr,
		"payload", payload,
	)
	return payload, true
}

// readPayload reads the body leniently: valid JSON is kept verbatim as
// json.RawMessage, anything else as a string, an empty body as nil.
// Tighten into a struct once the real Lumana payload is known.
func readPayload(w http.ResponseWriter, r *http.Request) (any, error) {
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxBodyBytes))
	switch {
	case err != nil:
		return nil, err
	case len(body) == 0:
		return nil, nil
	case json.Valid(body):
		return json.RawMessage(body), nil
	default:
		return string(body), nil
	}
}

func writeOK(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
