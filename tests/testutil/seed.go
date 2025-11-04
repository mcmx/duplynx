package testutil

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/mcmx/duplynx/ent"
	"github.com/mcmx/duplynx/ent/enttest"
	"github.com/mcmx/duplynx/internal/data"
)

// SeededClient bundles an Ent client with the canonical demo dataset.
type SeededClient struct {
	Client  *ent.Client
	Dataset data.DemoDataset
}

// NewSeededClient opens an in-memory Ent client, seeds the canonical dataset, and registers cleanup.
func NewSeededClient(t *testing.T) SeededClient {
	t.Helper()

	client := enttest.Open(t, "sqlite3", "file:duplynx-test?mode=memory&cache=shared&_fk=1")
	t.Cleanup(func() {
		_ = client.Close()
	})

	dataset := data.CanonicalDemoDataset()
	if _, err := data.SeedDemoDataset(context.Background(), client, dataset); err != nil {
		t.Fatalf("seed demo dataset: %v", err)
	}

	return SeededClient{
		Client:  client,
		Dataset: dataset,
	}
}

// TenantSlugFor returns the slug for the given tenant ID or fails the test if missing.
func TenantSlugFor(t *testing.T, dataset data.DemoDataset, tenantID uuid.UUID) string {
	t.Helper()
	for _, tenant := range dataset.Tenants {
		if tenant.ID == tenantID {
			return tenant.Slug
		}
	}
	t.Fatalf("tenant %s not found in dataset", tenantID)
	return ""
}

// MachineIDsForTenant returns all machine IDs for the tenant.
func MachineIDsForTenant(dataset data.DemoDataset, tenantID uuid.UUID) []uuid.UUID {
	ids := make([]uuid.UUID, 0)
	for _, machine := range dataset.Machines {
		if machine.TenantID == tenantID {
			ids = append(ids, machine.ID)
		}
	}
	return ids
}
