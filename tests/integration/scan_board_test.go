package integration_test

import (
	"context"
	"testing"

	"github.com/mcmx/duplynx/internal/scans"
)

func TestListTenantScans(t *testing.T) {
	repo := scans.NewRepository(scans.SampleScans())
	svc := scans.Service{ScansRepo: repo}
	items, err := svc.ListTenantScans(context.Background(), "sample-tenant-a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 3 {
		t.Fatalf("expected 3 scans, got %d", len(items))
	}
}
