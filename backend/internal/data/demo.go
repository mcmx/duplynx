package data

import (
	"time"

	"github.com/google/uuid"

	"github.com/mcmx/duplynx/internal/actions"
)

// DemoDataset captures the canonical demo records for the seed workflow.
type DemoDataset struct {
	Tenants         []TenantFixture
	Machines        []MachineFixture
	Scans           []ScanFixture
	DuplicateGroups []DuplicateGroupFixture
	FileInstances   []FileInstanceFixture
	ActionAudits    []ActionAuditFixture
}

// TenantFixture describes a tenant record in the canonical dataset.
type TenantFixture struct {
	ID             uuid.UUID
	Slug           string
	Name           string
	Description    string
	PrimaryContact string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// MachineFixture describes a machine record in the canonical dataset.
type MachineFixture struct {
	ID        uuid.UUID
	TenantID  uuid.UUID
	Name      string
	Category  string
	Hostname  string
	Role      string
	LastScan  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ScanFixture describes a scan record in the canonical dataset.
type ScanFixture struct {
	ID                  uuid.UUID
	TenantID            uuid.UUID
	InitiatedMachineID  uuid.UUID
	Name                string
	Description         string
	StartedAt           time.Time
	CompletedAt         time.Time
	DuplicateGroupCount int
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// DuplicateGroupFixture describes a duplicate group record in the canonical dataset.
type DuplicateGroupFixture struct {
	ID              uuid.UUID
	TenantID        uuid.UUID
	ScanID          uuid.UUID
	KeeperMachineID uuid.UUID
	Status          string
	Hash            string
	FileCount       int
	TotalSizeBytes  int64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// FileInstanceFixture describes file instances associated with duplicate groups.
type FileInstanceFixture struct {
	ID               uuid.UUID
	DuplicateGroupID uuid.UUID
	MachineID        uuid.UUID
	Path             string
	SizeBytes        int64
	Checksum         string
	LastSeenAt       time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Quarantined      bool
}

// ActionAuditFixture describes action audit records.
type ActionAuditFixture struct {
	ID               uuid.UUID
	TenantID         uuid.UUID
	DuplicateGroupID uuid.UUID
	Actor            string
	ActionType       string
	Payload          map[string]any
	PerformedAt      time.Time
	Stubbed          bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// CanonicalDemoDataset returns the canonical dataset for the demo workflow.
func CanonicalDemoDataset() DemoDataset {
	return canonicalDataset.Clone()
}

// Clone produces a deep copy of the dataset to avoid shared state during tests.
func (d DemoDataset) Clone() DemoDataset {
	clone := DemoDataset{
		Tenants:         append([]TenantFixture(nil), d.Tenants...),
		Machines:        append([]MachineFixture(nil), d.Machines...),
		Scans:           append([]ScanFixture(nil), d.Scans...),
		DuplicateGroups: append([]DuplicateGroupFixture(nil), d.DuplicateGroups...),
		FileInstances:   append([]FileInstanceFixture(nil), d.FileInstances...),
		ActionAudits:    append([]ActionAuditFixture(nil), d.ActionAudits...),
	}
	for i := range clone.ActionAudits {
		if d.ActionAudits[i].Payload != nil {
			payload := make(map[string]any, len(d.ActionAudits[i].Payload))
			for k, v := range d.ActionAudits[i].Payload {
				payload[k] = v
			}
			clone.ActionAudits[i].Payload = payload
		}
	}
	return clone
}

// DuplicateGroupsForStore converts duplicate group fixtures into in-memory action store records.
func (d DemoDataset) DuplicateGroupsForStore() []actions.DuplicateGroup {
	filesByGroup := make(map[uuid.UUID][]FileInstanceFixture)
	for _, file := range d.FileInstances {
		filesByGroup[file.DuplicateGroupID] = append(filesByGroup[file.DuplicateGroupID], file)
	}

	out := make([]actions.DuplicateGroup, 0, len(d.DuplicateGroups))
	for _, group := range d.DuplicateGroups {
		files := filesByGroup[group.ID]
		storeFiles := make([]actions.DuplicateFile, 0, len(files))
		for _, file := range files {
			storeFiles = append(storeFiles, actions.DuplicateFile{
				ID:          file.ID.String(),
				MachineID:   file.MachineID.String(),
				Path:        file.Path,
				SizeBytes:   file.SizeBytes,
				Quarantined: file.Quarantined,
			})
		}
		out = append(out, actions.DuplicateGroup{
			ID:              group.ID.String(),
			ScanID:          group.ScanID.String(),
			TenantSlug:      tenantSlugByID(group.TenantID, d.Tenants),
			Status:          group.Status,
			KeeperMachineID: group.KeeperMachineID.String(),
			Hash:            group.Hash,
			Files:           storeFiles,
		})
	}
	return out
}

func tenantSlugByID(id uuid.UUID, tenants []TenantFixture) string {
	for _, tenant := range tenants {
		if tenant.ID == id {
			return tenant.Slug
		}
	}
	return ""
}

var (
	tsOctober25 = time.Date(2025, 10, 25, 15, 30, 0, 0, time.UTC)
	tsOctober26 = time.Date(2025, 10, 26, 15, 30, 0, 0, time.UTC)
	tsOctober27 = time.Date(2025, 10, 27, 15, 30, 0, 0, time.UTC)
	tsOctober28 = time.Date(2025, 10, 28, 15, 30, 0, 0, time.UTC)

	tenantOrionID  = uuid.MustParse("5f2f6f52-1ac0-4a41-84b9-592effdb8b27")
	tenantSeleneID = uuid.MustParse("90d1b366-b2bb-4ed9-92ec-6c6ede98c24f")

	machineOrionCore    = uuid.MustParse("2b9a0101-7a92-4e81-b9a1-6c3a0d905f4a")
	machineOrionLaptop  = uuid.MustParse("a77ce4d0-1f02-44b5-bf23-8aea5c50d0c2")
	machineOrionArchive = uuid.MustParse("4cf02d63-574f-4f1b-b44c-6850b80b9a19")
	machineSeleneCore   = uuid.MustParse("1882ebba-9cb4-4f0d-9f52-17c8f7e54073")
	machineSeleneLab    = uuid.MustParse("2c6e9724-1b3d-47f8-8c5c-5e94eaaae5b0")

	scanBaseline   = uuid.MustParse("341971ac-a504-4cd0-af68-8b21f3796470")
	scanMediaAudit = uuid.MustParse("8e27c476-cccf-4b0a-979f-0a2a6c44f5de")
	scanDiscovery  = uuid.MustParse("1ec9b872-e31c-4644-829a-1602043614d0")

	groupFinanceDocs  = uuid.MustParse("c02a3a1b-50f1-4418-89fe-5983dce56250")
	groupMediaArchive = uuid.MustParse("e3b8fe2f-a6da-407d-96fd-1aea6a0c97c4")
	groupDesignShare  = uuid.MustParse("0bde87e4-0170-4de4-b63c-6281a65c1dc9")
	groupTelemetry    = uuid.MustParse("4c9581d5-6ed4-4820-9e6e-7a6a5f239f78")

	fileFinanceLaptop  = uuid.MustParse("7bc9e001-d2de-4ed4-9424-9489d38c8334")
	fileFinanceShare   = uuid.MustParse("b2bac5f8-a6b1-4f9d-bb3a-9cdd3d67b6dd")
	fileFinanceArchive = uuid.MustParse("a5cfa963-7439-4bb8-9eb2-78d0871c4692")

	fileMediaOne     = uuid.MustParse("d4da6d10-3fe5-4a73-9818-7ea85f992e4b")
	fileMediaTwo     = uuid.MustParse("e8330a58-3b39-4f1f-931b-4bcdcd8e0f09")
	fileDesignUX     = uuid.MustParse("ccdb5a06-5b83-4cf0-803c-96e4034f068f")
	fileDesignDeck   = uuid.MustParse("a8422d51-583a-4d3f-bb86-f5f7cc149048")
	fileTelemetry    = uuid.MustParse("0f9a2c09-2fd9-4a8f-9c5a-0f7f43c6b3d4")
	fileTelemetryLab = uuid.MustParse("037620d8-0ed6-4a8e-9214-0655e768f194")

	auditAssignKeeper = uuid.MustParse("11ea57df-9a16-46f6-834e-0785d08c0f3b")
	auditQuarantine   = uuid.MustParse("320c6e29-c9df-4f3e-a2a9-5688b1ec6764")
)

var canonicalDataset = DemoDataset{
	Tenants: []TenantFixture{
		{
			ID:             tenantOrionID,
			Slug:           "orion-analytics",
			Name:           "Orion Analytics",
			Description:    "Northwind subsidiary focused on data deduplication rollouts",
			PrimaryContact: "demo@orion.test",
			CreatedAt:      tsOctober25,
			UpdatedAt:      tsOctober28,
		},
		{
			ID:             tenantSeleneID,
			Slug:           "selene-research",
			Name:           "Selene Research",
			Description:    "Applied sciences lab maintaining mixed media archives",
			PrimaryContact: "ops@selene.test",
			CreatedAt:      tsOctober25,
			UpdatedAt:      tsOctober27,
		},
	},
	Machines: []MachineFixture{
		{
			ID:        machineOrionCore,
			TenantID:  tenantOrionID,
			Name:      "Orion Core 01",
			Category:  "server",
			Hostname:  "orion-core-01.orion.test",
			Role:      "ingest",
			LastScan:  time.Date(2025, 10, 28, 9, 30, 0, 0, time.UTC),
			CreatedAt: tsOctober25,
			UpdatedAt: tsOctober28,
		},
		{
			ID:        machineOrionLaptop,
			TenantID:  tenantOrionID,
			Name:      "Orion Laptop 01",
			Category:  "personal_laptop",
			Hostname:  "laptop-01.orion.test",
			Role:      "analysis",
			LastScan:  time.Date(2025, 10, 27, 17, 5, 0, 0, time.UTC),
			CreatedAt: tsOctober25,
			UpdatedAt: tsOctober27,
		},
		{
			ID:        machineOrionArchive,
			TenantID:  tenantOrionID,
			Name:      "Orion Archive 01",
			Category:  "server",
			Hostname:  "archive-01.orion.test",
			Role:      "archive",
			LastScan:  time.Date(2025, 10, 28, 7, 55, 0, 0, time.UTC),
			CreatedAt: tsOctober25,
			UpdatedAt: tsOctober28,
		},
		{
			ID:        machineSeleneCore,
			TenantID:  tenantSeleneID,
			Name:      "Selene Core 01",
			Category:  "server",
			Hostname:  "core-01.selene.test",
			Role:      "render",
			LastScan:  time.Date(2025, 10, 27, 10, 40, 0, 0, time.UTC),
			CreatedAt: tsOctober25,
			UpdatedAt: tsOctober27,
		},
		{
			ID:        machineSeleneLab,
			TenantID:  tenantSeleneID,
			Name:      "Selene Lab 01",
			Category:  "personal_laptop",
			Hostname:  "lab-01.selene.test",
			Role:      "fieldwork",
			LastScan:  time.Date(2025, 10, 26, 20, 5, 0, 0, time.UTC),
			CreatedAt: tsOctober25,
			UpdatedAt: tsOctober26,
		},
	},
	Scans: []ScanFixture{
		{
			ID:                  scanBaseline,
			TenantID:            tenantOrionID,
			InitiatedMachineID:  machineOrionCore,
			Name:                "Baseline Sweep October",
			Description:         "Weekly deduplication sweep across finance shares",
			StartedAt:           time.Date(2025, 10, 28, 8, 0, 0, 0, time.UTC),
			CompletedAt:         time.Date(2025, 10, 28, 8, 22, 0, 0, time.UTC),
			DuplicateGroupCount: 2,
			CreatedAt:           tsOctober28,
			UpdatedAt:           tsOctober28,
		},
		{
			ID:                  scanMediaAudit,
			TenantID:            tenantOrionID,
			InitiatedMachineID:  machineOrionArchive,
			Name:                "Media Archive Audit",
			Description:         "Quarterly media sync validation",
			StartedAt:           time.Date(2025, 10, 26, 22, 5, 0, 0, time.UTC),
			CompletedAt:         time.Date(2025, 10, 26, 22, 47, 0, 0, time.UTC),
			DuplicateGroupCount: 1,
			CreatedAt:           tsOctober26,
			UpdatedAt:           tsOctober26,
		},
		{
			ID:                  scanDiscovery,
			TenantID:            tenantSeleneID,
			InitiatedMachineID:  machineSeleneCore,
			Name:                "Telemetry Discovery",
			Description:         "Nightly telemetry shape check for field devices",
			StartedAt:           time.Date(2025, 10, 27, 10, 15, 0, 0, time.UTC),
			CompletedAt:         time.Date(2025, 10, 27, 10, 38, 0, 0, time.UTC),
			DuplicateGroupCount: 1,
			CreatedAt:           tsOctober27,
			UpdatedAt:           tsOctober27,
		},
	},
	DuplicateGroups: []DuplicateGroupFixture{
		{
			ID:              groupFinanceDocs,
			TenantID:        tenantOrionID,
			ScanID:          scanBaseline,
			KeeperMachineID: machineOrionCore,
			Status:          "review",
			Hash:            "hash:finance-q4-plan",
			FileCount:       3,
			TotalSizeBytes:  3_145_728,
			CreatedAt:       tsOctober28,
			UpdatedAt:       tsOctober28,
		},
		{
			ID:              groupMediaArchive,
			TenantID:        tenantOrionID,
			ScanID:          scanBaseline,
			KeeperMachineID: machineOrionArchive,
			Status:          "action_needed",
			Hash:            "hash:media-latest-cut",
			FileCount:       2,
			TotalSizeBytes:  8_388_608,
			CreatedAt:       tsOctober28,
			UpdatedAt:       tsOctober28,
		},
		{
			ID:              groupDesignShare,
			TenantID:        tenantOrionID,
			ScanID:          scanMediaAudit,
			KeeperMachineID: machineOrionLaptop,
			Status:          "resolved",
			Hash:            "hash:design-review-pack",
			FileCount:       2,
			TotalSizeBytes:  2_621_440,
			CreatedAt:       tsOctober26,
			UpdatedAt:       tsOctober26,
		},
		{
			ID:              groupTelemetry,
			TenantID:        tenantSeleneID,
			ScanID:          scanDiscovery,
			KeeperMachineID: machineSeleneCore,
			Status:          "review",
			Hash:            "hash:telemetry-drift",
			FileCount:       2,
			TotalSizeBytes:  1_572_864,
			CreatedAt:       tsOctober27,
			UpdatedAt:       tsOctober27,
		},
	},
	FileInstances: []FileInstanceFixture{
		{
			ID:               fileFinanceLaptop,
			DuplicateGroupID: groupFinanceDocs,
			MachineID:        machineOrionLaptop,
			Path:             "/Users/finance/roadmap/Q4-plan.pptx",
			SizeBytes:        1_048_576,
			Checksum:         "sha256:9adcc6e5a3b9f83bbd1c52c4182f92a5",
			LastSeenAt:       time.Date(2025, 10, 28, 8, 15, 0, 0, time.UTC),
			CreatedAt:        tsOctober28,
			UpdatedAt:        tsOctober28,
		},
		{
			ID:               fileFinanceShare,
			DuplicateGroupID: groupFinanceDocs,
			MachineID:        machineOrionCore,
			Path:             "/srv/shares/finance/Q4-plan.pptx",
			SizeBytes:        1_048_576,
			Checksum:         "sha256:9adcc6e5a3b9f83bbd1c52c4182f92a5",
			LastSeenAt:       time.Date(2025, 10, 28, 8, 16, 0, 0, time.UTC),
			CreatedAt:        tsOctober28,
			UpdatedAt:        tsOctober28,
		},
		{
			ID:               fileFinanceArchive,
			DuplicateGroupID: groupFinanceDocs,
			MachineID:        machineOrionArchive,
			Path:             "/srv/archive/finance/Q4-plan.pptx",
			SizeBytes:        1_048_576,
			Checksum:         "sha256:9adcc6e5a3b9f83bbd1c52c4182f92a5",
			LastSeenAt:       time.Date(2025, 10, 28, 8, 16, 30, 0, time.UTC),
			CreatedAt:        tsOctober28,
			UpdatedAt:        tsOctober28,
		},
		{
			ID:               fileMediaOne,
			DuplicateGroupID: groupMediaArchive,
			MachineID:        machineOrionArchive,
			Path:             "/srv/media/cuts/fall-launch.mov",
			SizeBytes:        4_194_304,
			Checksum:         "sha256:be93b99045db02b9301a51597d90f4fb",
			LastSeenAt:       time.Date(2025, 10, 28, 7, 30, 0, 0, time.UTC),
			CreatedAt:        tsOctober28,
			UpdatedAt:        tsOctober28,
		},
		{
			ID:               fileMediaTwo,
			DuplicateGroupID: groupMediaArchive,
			MachineID:        machineOrionCore,
			Path:             "/srv/media/staging/fall-launch.mov",
			SizeBytes:        4_194_304,
			Checksum:         "sha256:be93b99045db02b9301a51597d90f4fb",
			LastSeenAt:       time.Date(2025, 10, 28, 7, 31, 0, 0, time.UTC),
			CreatedAt:        tsOctober28,
			UpdatedAt:        tsOctober28,
			Quarantined:      true,
		},
		{
			ID:               fileDesignUX,
			DuplicateGroupID: groupDesignShare,
			MachineID:        machineOrionLaptop,
			Path:             "/Users/design/review/package.zip",
			SizeBytes:        1_310_720,
			Checksum:         "sha256:e845d474d887cf43d62f56d8f762b73f",
			LastSeenAt:       time.Date(2025, 10, 26, 22, 20, 0, 0, time.UTC),
			CreatedAt:        tsOctober26,
			UpdatedAt:        tsOctober26,
			Quarantined:      true,
		},
		{
			ID:               fileDesignDeck,
			DuplicateGroupID: groupDesignShare,
			MachineID:        machineOrionArchive,
			Path:             "/srv/design/latest/package.zip",
			SizeBytes:        1_310_720,
			Checksum:         "sha256:e845d474d887cf43d62f56d8f762b73f",
			LastSeenAt:       time.Date(2025, 10, 26, 22, 25, 0, 0, time.UTC),
			CreatedAt:        tsOctober26,
			UpdatedAt:        tsOctober26,
		},
		{
			ID:               fileTelemetry,
			DuplicateGroupID: groupTelemetry,
			MachineID:        machineSeleneCore,
			Path:             "/srv/telemetry/batch-219.json",
			SizeBytes:        786_432,
			Checksum:         "sha256:49ba5eb796f5a6d0750f0dcd384e0b8e",
			LastSeenAt:       time.Date(2025, 10, 27, 10, 30, 0, 0, time.UTC),
			CreatedAt:        tsOctober27,
			UpdatedAt:        tsOctober27,
		},
		{
			ID:               fileTelemetryLab,
			DuplicateGroupID: groupTelemetry,
			MachineID:        machineSeleneLab,
			Path:             "/Users/lab/cache/batch-219.json",
			SizeBytes:        786_432,
			Checksum:         "sha256:49ba5eb796f5a6d0750f0dcd384e0b8e",
			LastSeenAt:       time.Date(2025, 10, 27, 10, 32, 0, 0, time.UTC),
			CreatedAt:        tsOctober27,
			UpdatedAt:        tsOctober27,
		},
	},
	ActionAudits: []ActionAuditFixture{
		{
			ID:               auditAssignKeeper,
			TenantID:         tenantOrionID,
			DuplicateGroupID: groupFinanceDocs,
			Actor:            "demo",
			ActionType:       "assign_keeper",
			Payload: map[string]any{
				"keeperMachineId": machineOrionCore.String(),
			},
			PerformedAt: time.Date(2025, 10, 28, 8, 30, 0, 0, time.UTC),
			Stubbed:     false,
			CreatedAt:   tsOctober28,
			UpdatedAt:   tsOctober28,
		},
		{
			ID:               auditQuarantine,
			TenantID:         tenantOrionID,
			DuplicateGroupID: groupMediaArchive,
			Actor:            "demo",
			ActionType:       "quarantine",
			Payload: map[string]any{
				"reason": "media sync window overlap",
			},
			PerformedAt: time.Date(2025, 10, 28, 7, 45, 0, 0, time.UTC),
			Stubbed:     true,
			CreatedAt:   tsOctober28,
			UpdatedAt:   tsOctober28,
		},
	},
}
