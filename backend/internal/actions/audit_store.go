package actions

// AuditStore persists audit entries (in addition to AuditLogger slice) for demo retrieval.
type AuditStore struct {
	Logger *AuditLogger
}

func (s AuditStore) Append(entry AuditEntry) {
	if s.Logger != nil {
		s.Logger.Log(entry)
	}
}

func (s AuditStore) Entries() []AuditEntry {
	if s.Logger == nil {
		return nil
	}
	return s.Logger.Entries()
}
