package integration_test

import (
	"context"
	"testing"

	"github.com/mcmx/duplynx/internal/tenancy"
)

func TestListMachinesFiltersByTenant(t *testing.T) {
	repo := tenancy.NewRepository(tenancy.SampleTenants(), nil)
	machines, err := repo.ListMachines(context.Background(), "sample-tenant-a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(machines) != 5 {
		t.Fatalf("expected 5 machines, got %d", len(machines))
	}
}
