package smoke_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestSmokeSeedServe(t *testing.T) {
	repoRoot := findRepoRoot(t)
	backendDir := filepath.Join(repoRoot, "backend")

	tempDir := t.TempDir()
	dbFile := filepath.Join(tempDir, "duplynx_smoke.db")
	assetsDir := filepath.Join(tempDir, "assets")
	if err := os.MkdirAll(assetsDir, 0o755); err != nil {
		t.Fatalf("failed to create assets dir: %v", err)
	}
	cssContent := []byte("/* smoke asset */")
	if err := os.WriteFile(filepath.Join(assetsDir, "tailwind.css"), cssContent, 0o644); err != nil {
		t.Fatalf("failed to scaffold tailwind bundle: %v", err)
	}
	if err := os.WriteFile(filepath.Join(assetsDir, "app.css"), cssContent, 0o644); err != nil {
		t.Fatalf("failed to scaffold app.css bundle: %v", err)
	}

	env := append(os.Environ(),
		"GOMODCACHE="+filepath.Join(repoRoot, ".cache", "go-mod"),
		"GOCACHE="+filepath.Join(repoRoot, ".cache", "go-build"),
		"GOSUMDB=off",
	)

	seedStart := time.Now()
	runCommand(t, backendDir, env, time.Minute, "go", "run", "./cmd/duplynx", "seed",
		"--db-file", dbFile,
		"--assets-dir", assetsDir,
	)
	seedDuration := time.Since(seedStart)
	if seedDuration > time.Minute {
		t.Fatalf("seed command exceeded 60s budget: %s", seedDuration)
	}

	addr := fmt.Sprintf("127.0.0.1:%d", freePort(t))
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	var stdout bytes.Buffer
	cmd := exec.CommandContext(ctx, "go", "run", "./cmd/duplynx", "serve",
		"--db-file", dbFile,
		"--assets-dir", assetsDir,
		"--addr", addr,
	)
	if runtime.GOOS != "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	}
	cmd.Dir = backendDir
	cmd.Env = env
	cmd.Stdout = &stdout
	cmd.Stderr = &stdout

	serveStart := time.Now()
	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start serve command: %v", err)
	}

	waitCtx, waitCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer waitCancel()
	if err := waitForServer(waitCtx, "http://"+addr+"/healthz"); err != nil {
		_ = cmd.Process.Kill()
		t.Fatalf("server did not become ready: %v\noutput:\n%s", err, stdout.String())
	}

	rootResp, err := http.Get("http://" + addr + "/")
	if err != nil {
		_ = cmd.Process.Signal(os.Interrupt)
		t.Fatalf("failed to call root route: %v", err)
	}
	defer rootResp.Body.Close()

	body, err := io.ReadAll(rootResp.Body)
	if err != nil {
		_ = cmd.Process.Signal(os.Interrupt)
		t.Fatalf("failed to read root response: %v", err)
	}
	if rootResp.StatusCode != http.StatusOK {
		_ = cmd.Process.Signal(os.Interrupt)
		t.Fatalf("expected 200 status, got %d.\nOutput:\n%s", rootResp.StatusCode, stdout.String())
	}

	html := string(body)
	if !strings.Contains(html, "/static/app.css") {
		_ = cmd.Process.Signal(os.Interrupt)
		t.Fatalf("expected Tailwind bundle link in markup.\nReceived HTML:\n%s", truncate(html, 512))
	}
	if !strings.Contains(html, "Orion Analytics") {
		t.Fatalf("expected seeded tenant name in markup.\nReceived HTML:\n%s", truncate(html, 512))
	}

	if err := sendGroupSignal(cmd, syscall.SIGINT); err != nil && !errors.Is(err, os.ErrProcessDone) {
		t.Fatalf("failed to deliver SIGINT to serve process: %v", err)
	}
	// Provide a follow-up SIGTERM to ensure shutdown on platforms that ignore SIGINT.
	_ = sendGroupSignal(cmd, syscall.SIGTERM)
	cancel()

	waitCh := make(chan error, 1)
	go func() {
		waitCh <- cmd.Wait()
	}()

	select {
	case err := <-waitCh:
		if err != nil && !isBenignExit(err) {
			t.Fatalf("serve command exited with error: %v\noutput:\n%s", err, stdout.String())
		}
	case <-time.After(10 * time.Second):
		_ = killProcessGroup(cmd)
		if err := <-waitCh; err != nil && !isBenignExit(err) {
			t.Fatalf("serve command did not exit cleanly after SIGKILL: %v\noutput:\n%s", err, stdout.String())
		}
		t.Log("serve command required SIGKILL to terminate within smoke timeout")
	}

	serverDuration := time.Since(serveStart)
	if serverDuration > 2*time.Minute {
		t.Fatalf("serve lifecycle exceeded 2 minute budget: %s", serverDuration)
	}
}

func runCommand(t *testing.T, dir string, env []string, timeout time.Duration, name string, args ...string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	cmd.Env = env
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command %s %v failed: %v\n%s", name, args, err, string(output))
	}
}

func findRepoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to resolve working directory: %v", err)
	}
	// tests module root is .../tests/smoke
	for {
		if dir == "/" || dir == "" {
			t.Fatalf("unable to locate repository root from %s", dir)
		}
		if _, err := os.Stat(filepath.Join(dir, "Makefile")); err == nil {
			return dir
		}
		dir = filepath.Dir(dir)
	}
}

func freePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to allocate port: %v", err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

func waitForServer(ctx context.Context, url string) error {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	client := http.Client{Timeout: 2 * time.Second}
	for {
		resp, err := client.Get(url)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("context ended waiting for %s: %w", url, ctx.Err())
		case <-ticker.C:
		}
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "…"
}

func isBenignExit(err error) bool {
	if err == nil {
		return true
	}
	if errors.Is(err, context.Canceled) {
		return true
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		// Unix negative values represent signals; treat signal-triggered shutdown as success.
		return exitErr.ExitCode() == 0 || exitErr.ExitCode() == -1
	}
	return false
}

func sendGroupSignal(cmd *exec.Cmd, sig syscall.Signal) error {
	if cmd.Process == nil {
		return errors.New("process not started")
	}
	if runtime.GOOS == "windows" {
		return cmd.Process.Signal(sig)
	}
	return syscall.Kill(-cmd.Process.Pid, sig)
}

func killProcessGroup(cmd *exec.Cmd) error {
	if cmd.Process == nil {
		return errors.New("process not started")
	}
	if runtime.GOOS == "windows" {
		return cmd.Process.Kill()
	}
	return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
}
