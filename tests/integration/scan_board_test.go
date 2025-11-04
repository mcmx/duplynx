package integration_test

import (
	"context"
	"testing"

	"github.com/mcmx/duplynx/internal/scans"
	"github.com/mcmx/duplynx/tests/testutil"
)

func TestListTenantScans(t *testing.T) {
	seed := testutil.NewSeededClient(t)
	repo := scans.NewRepositoryFromClient(seed.Client)
	svc := scans.Service{Repo: repo}
	tenantSlug := seed.Dataset.Tenants[0].Slug
	items, err := svc.ListTenantScans(context.Background(), tenantSlug)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) == 0 {
		t.Fatalf("expected scans for tenant %s", tenantSlug)
	}
}
