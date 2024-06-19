package slogdiscard

import (
	"context"
	"log/slog"
)

// NewDiscardLogger returns a logger that discards all log entries
func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

// DiscardHandler is a handler that discards all log entries
type DiscardHandler struct{}

// NewDiscardHandler returns a new DiscardHandler
func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}

// Handle just ignores the log entry
func (h *DiscardHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

// WithAttrs returns the sane handler since there are no attributes to save
func (h *DiscardHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

// WithGroup returns the same handler since there is no group to save
func (h *DiscardHandler) WithGroup(_ string) slog.Handler {
	return h
}

// Enabled returns false since we don't want to log anything
func (h *DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}
