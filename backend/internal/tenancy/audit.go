package tenancy

import "sync"

// AuditEntry captures a user-triggered tenant or machine selection event.
type AuditEntry struct {
	Type        string
	TenantSlug  string
	MachineID   string
	MachineName string
}

// AuditLogger stores audit entries in-memory for demo purposes.
type AuditLogger struct {
	mu      sync.Mutex
	entries []AuditEntry
}

func (a *AuditLogger) LogTenantSelection(slug string) {
	if a == nil {
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.entries = append(a.entries, AuditEntry{Type: "tenant_selection", TenantSlug: slug})
}

func (a *AuditLogger) LogMachineSelection(slug string, machine Machine) {
	if a == nil {
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.entries = append(a.entries, AuditEntry{
		Type:        "machine_selection",
		TenantSlug:  slug,
		MachineID:   machine.ID,
		MachineName: machine.Name,
	})
}

// Entries returns a copy of the audit entries slice for assertions.
func (a *AuditLogger) Entries() []AuditEntry {
	if a == nil {
		return nil
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	out := make([]AuditEntry, len(a.entries))
	copy(out, a.entries)
	return out
}
