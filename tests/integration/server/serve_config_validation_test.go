package integration_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestServeFailsWhenAssetsMissing(t *testing.T) {
	repoRoot, err := filepath.Abs(".")
	if err != nil {
		t.Fatalf("failed to determine repo root: %v", err)
	}
	backendDir := filepath.Join(repoRoot, "backend")
	tempDir := t.TempDir()
	missingAssets := filepath.Join(tempDir, "no-assets")
	if err := os.MkdirAll(missingAssets, 0o755); err != nil {
		t.Fatalf("failed to create temp asset dir: %v", err)
	}
	dbPath := filepath.Join(tempDir, "duplynx.db")

	cmd := exec.Command("go", "run", "./cmd/duplynx", "serve", "--assets-dir", missingAssets, "--db-file", dbPath)
	cmd.Dir = backendDir
	cmd.Env = append(os.Environ(),
		"GOMODCACHE="+filepath.Join(repoRoot, ".cache", "go-mod"),
		"GOCACHE="+filepath.Join(repoRoot, ".cache", "go-build"),
		"GOSUMDB=off",
	)

	output, runErr := cmd.CombinedOutput()
	if runErr == nil {
		t.Fatalf("expected serve command to fail when assets are missing, output: %s", string(output))
	}

	if !containsString(string(output), "tailwind bundle missing") {
		t.Fatalf("expected missing asset error, got: %s", string(output))
	}
}

func containsString(haystack, needle string) bool {
	return strings.Contains(haystack, needle)
}
