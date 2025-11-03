package integration_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	entmachine "github.com/mcmx/duplynx/ent/machine"
	enttenant "github.com/mcmx/duplynx/ent/tenant"
	"github.com/mcmx/duplynx/internal/config"
	"github.com/mcmx/duplynx/internal/data"
)

func TestSeedCommandPopulatesCanonicalDataset(t *testing.T) {
	repoRoot, err := filepath.Abs(".")
	if err != nil {
		t.Fatalf("failed to determine repo root: %v", err)
	}

	backendDir := filepath.Join(repoRoot, "backend")
	assetsDir := filepath.Join(backendDir, "web", "dist")
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "duplynx.db")

	runSeedCommand(t, backendDir, assetsDir, dbPath)
	runSeedCommand(t, backendDir, assetsDir, dbPath)

	ctx := context.Background()
	cfg := config.RuntimeConfig{DBFile: dbPath}

	client, err := data.OpenSQLite(ctx, cfg.SQLiteDSN())
	if err != nil {
		t.Fatalf("failed to open sqlite database: %v", err)
	}
	defer func() {
		_ = data.Close(client)
	}()

	dataset := data.CanonicalDemoDataset()

	tenants, err := client.Tenant.Query().Order(enttenant.BySlug()).All(ctx)
	if err != nil {
		t.Fatalf("failed to query tenants: %v", err)
	}

	if len(tenants) != len(dataset.Tenants) {
		t.Fatalf("expected %d tenants, got %d", len(dataset.Tenants), len(tenants))
	}

	expectedMachines := make(map[string]int)
	for _, machine := range dataset.Machines {
		expectedMachines[machine.TenantID.String()]++
	}

	for _, tenantRecord := range tenants {
		expectedCount := expectedMachines[tenantRecord.ID.String()]
		count, err := client.Machine.
			Query().
			Where(entmachine.HasTenantWith(enttenant.IDEQ(tenantRecord.ID))).
			Count(ctx)
		if err != nil {
			t.Fatalf("failed to count machines for tenant %s: %v", tenantRecord.Slug, err)
		}
		if count != expectedCount {
			t.Fatalf("expected %d machines for tenant %s, got %d", expectedCount, tenantRecord.Slug, count)
		}
	}
}

func runSeedCommand(t *testing.T, backendDir, assetsDir, dbPath string) {
	t.Helper()

	cmd := exec.Command("go", "run", "./cmd/duplynx", "seed", "--db-file", dbPath, "--assets-dir", assetsDir)
	cmd.Dir = backendDir
	cmd.Env = append(os.Environ(),
		"GOMODCACHE="+filepath.Join(backendDir, "..", ".cache", "go-mod"),
		"GOCACHE="+filepath.Join(backendDir, "..", ".cache", "go-build"),
		"GOSUMDB=off",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("seed command failed: %v\n%s", err, string(output))
	}
}
