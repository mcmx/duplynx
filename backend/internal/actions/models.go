package actions

// DuplicateGroup represents a collection of duplicate files detected by a scan.
type DuplicateGroup struct {
	ID              string
	ScanID          string
	TenantSlug      string
	Status          string
	KeeperMachineID string
	Hash            string
	Files           []DuplicateFile
}

// DuplicateFile describes an instance of the duplicate within a machine.
type DuplicateFile struct {
	ID          string
	MachineID   string
	Path        string
	SizeBytes   int64
	Quarantined bool
}
