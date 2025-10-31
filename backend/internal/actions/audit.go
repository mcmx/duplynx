package actions

import "sync"

// AuditEntry captures details about keeper assignments and duplicate actions.
type AuditEntry struct {
	Type            string
	GroupID         string
	TenantSlug      string
	Actor           string
	KeeperMachineID string
	Payload         map[string]any
	Stubbed         bool
}

// AuditLogger records audit entries in memory for demo purposes.
type AuditLogger struct {
	mu      sync.Mutex
	entries []AuditEntry
}

func (l *AuditLogger) Log(entry AuditEntry) {
	if l == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = append(l.entries, entry)
}

func (l *AuditLogger) Entries() []AuditEntry {
	if l == nil {
		return nil
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]AuditEntry, len(l.entries))
	copy(out, l.entries)
	return out
}
