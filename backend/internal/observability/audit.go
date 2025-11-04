package observability

import (
	"log"
	"log/slog"
	"os"
	"time"
)

// Event represents a lifecycle action (seed, serve, etc.) recorded by the CLI.
type Event struct {
	Action   string
	Actor    string
	Outcome  string
	Duration time.Duration
	Metadata map[string]any
	Error    error
}

// EventWriter records audit events to stdout via slog for consistent observability.
type EventWriter struct {
	logger *slog.Logger
}

// NewEventWriter constructs a writer using the supplied slog logger or a default text logger.
func NewEventWriter(logger *slog.Logger) EventWriter {
	if logger == nil {
		handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: false})
		logger = slog.New(handler)
	}
	return EventWriter{logger: logger}
}

// Start begins timing an action and returns a scope that records the outcome when finished.
func (w EventWriter) Start(action, actor string, metadata map[string]any) *Scope {
	if metadata == nil {
		metadata = map[string]any{}
	}
	return &Scope{
		writer:   w,
		action:   action,
		actor:    actor,
		started:  time.Now(),
		metadata: metadata,
	}
}

// Write emits a one-off event without timing helpers.
func (w EventWriter) Write(event Event) {
	if event.Metadata == nil {
		event.Metadata = map[string]any{}
	}

	if w.logger == nil {
		log.Printf("audit event %s actor=%s outcome=%s duration=%s metadata=%v err=%v",
			event.Action, event.Actor, event.Outcome, event.Duration, event.Metadata, event.Error)
		return
	}

	args := []any{
		slog.String("action", event.Action),
		slog.String("actor", event.Actor),
		slog.String("outcome", event.Outcome),
		slog.Duration("duration", event.Duration),
	}
	if event.Error != nil {
		args = append(args, slog.String("error", event.Error.Error()))
	}
	if len(event.Metadata) > 0 {
		args = append(args, slog.Any("metadata", event.Metadata))
	}
	w.logger.Info("audit event", args...)
}

// Scope captures timing for a long-running action and writes an audit entry when finished.
type Scope struct {
	writer   EventWriter
	action   string
	actor    string
	started  time.Time
	metadata map[string]any
}

// Finish records the outcome of the scoped action.
func (s *Scope) Finish(outcome string, err error) {
	if s == nil {
		return
	}
	event := Event{
		Action:   s.action,
		Actor:    s.actor,
		Outcome:  outcome,
		Duration: time.Since(s.started),
		Metadata: s.metadata,
		Error:    err,
	}
	s.writer.Write(event)
}
