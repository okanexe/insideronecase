// Package api provides HTTP handlers for message sending and management.
package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

// StartMessageSending godoc
// @Summary Start message sending worker
// @Description Starts the background worker that sends messages
// @Success 200 {string} string "Message sending started"
// @Failure 500 {string} string
// @Router /start [post]
func (a *API) StartMessageSending(w http.ResponseWriter, r *http.Request) {
	err := a.s.Start()
	if err != nil {
		slog.Error("failed to start message sending", slog.Any("error", err))
		http.Error(w, fmt.Sprintf("failed to start message sending: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Message sending started"))
	if err != nil {
		slog.Error("failed to write start message sending", slog.Any("error", err))
		http.Error(w, fmt.Sprintf("failed to write start message sending: %v", err), http.StatusInternalServerError)
		return
	}
}

// StopMessageSending godoc
// @Summary Stop message sending worker
// @Description Stops the background worker
// @Success 200 {string} string "Message sending stopped"
// @Failure 500 {string} string
// @Router /stop [post]
func (a *API) StopMessageSending(w http.ResponseWriter, r *http.Request) {
	err := a.s.Stop()
	if err != nil {
		slog.Error("failed to stop message sending", slog.Any("error", err))
		http.Error(w, fmt.Sprintf("failed to stop message sending: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Message sending stopped"))
	if err != nil {
		slog.Error("failed to write stop message sending", slog.Any("error", err))
		http.Error(w, fmt.Sprintf("failed to write stop message sending: %v", err), http.StatusInternalServerError)
		return
	}
}

// GetSentMessages godoc
// @Summary Get sent messages
// @Description Returns list of sent messages
// @Success 200 {array} cache.SentMessage
// @Failure 500 {string} string
// @Router /messages [get]
func (a *API) GetSentMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	m, err := a.s.Cache.GetMessages(ctx)
	if err != nil {
		slog.Error("failed to retrieve sent messages", slog.Any("error", err))
		http.Error(w, "Failed to retrieve sent messages", http.StatusInternalServerError)
		return
	}

	if len(m) == 0 {
		slog.Error("sent messages returned an empty sent messages")
		http.Error(w, fmt.Sprintf("No sent messages found"), http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(m)
	if err != nil {
		slog.Error("failed to encode get sent messages", slog.Any("error", err))
		http.Error(w, fmt.Sprintf("failed to encode get sent messages: %v", err), http.StatusInternalServerError)
		return
	}
}
