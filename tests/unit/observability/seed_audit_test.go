package observability_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/mcmx/duplynx/internal/observability"
)

func TestSeedAuditEventsIncludeActorAndOutcome(t *testing.T) {
	handler := &recordingHandler{}
	writer := observability.NewEventWriter(slog.New(handler))

	writer.Write(observability.Event{
		Action:  "seed_start",
		Actor:   "tester",
		Outcome: "starting",
	})
	writer.Write(observability.Event{
		Action:  "seed_stop",
		Actor:   "tester",
		Outcome: "success",
	})

	if len(handler.records) != 2 {
		t.Fatalf("expected 2 audit records, got %d", len(handler.records))
	}

	for _, record := range handler.records {
		action := attributeValue(record, "action")
		actor := attributeValue(record, "actor")
		outcome := attributeValue(record, "outcome")

		if action == "" {
			t.Fatalf("audit record missing action attribute: %+v", record)
		}
		if actor == "" {
			t.Fatalf("audit record missing actor attribute: %+v", record)
		}
		if outcome == "" {
			t.Fatalf("audit record missing outcome attribute: %+v", record)
		}
	}
}

type recordingHandler struct {
	records []slog.Record
}

func (h *recordingHandler) Enabled(context.Context, slog.Level) bool { return true }

func (h *recordingHandler) Handle(_ context.Context, record slog.Record) error {
	h.records = append(h.records, record.Clone())
	return nil
}

func (h *recordingHandler) WithAttrs([]slog.Attr) slog.Handler {
	return h
}

func (h *recordingHandler) WithGroup(string) slog.Handler {
	return h
}

func attributeValue(record slog.Record, key string) string {
	var value string
	record.Attrs(func(attr slog.Attr) bool {
		if attr.Key == key {
			value = attr.Value.String()
			return false
		}
		return true
	})
	return value
}
