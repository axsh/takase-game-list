package integration

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildOnlyMode_Integration(t *testing.T) {
	cleanBin(t)

	output, err := runBuildScript(t, "-Mode", "build-only")
	if err != nil {
		t.Fatalf("build script failed (build-only): %v\noutput:\n%s", err, output)
	}

	binPath := filepath.Join(projectRoot(t), "bin", "game-list.exe")
	if _, err := os.Stat(binPath); err != nil {
		t.Fatalf("expected binary at %s: %v\noutput:\n%s", binPath, err, output)
	}

	if !strings.Contains(output, "Mode: build-only") {
		t.Fatalf("expected output to include mode info, got:\n%s", output)
	}
	if !strings.Contains(output, "Build succeeded") {
		t.Fatalf("expected build success log, got:\n%s", output)
	}
}

func TestDefaultAllMode_Integration(t *testing.T) {
	cleanBin(t)

	output, err := runBuildScript(t)
	if err != nil {
		t.Fatalf("build script failed (default all): %v\noutput:\n%s", err, output)
	}

	if !strings.Contains(output, "Mode: all") {
		t.Fatalf("expected mode=all log, got:\n%s", output)
	}
	if !strings.Contains(output, "Running unit tests") {
		t.Fatalf("expected unit test log, got:\n%s", output)
	}
	if !strings.Contains(output, "Running integration tests") {
		t.Fatalf("expected integration test log, got:\n%s", output)
	}
}

func runBuildScript(t *testing.T, args ...string) (string, error) {
	t.Helper()

	cmdArgs := append([]string{
		"-NoProfile",
		"-ExecutionPolicy", "Bypass",
		"-File", "build.ps1",
	}, args...)

	cmd := exec.Command("powershell", cmdArgs...)
	cmd.Dir = projectRoot(t)

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	err := cmd.Run()
	return buf.String(), err
}

func cleanBin(t *testing.T) {
	t.Helper()
	binDir := filepath.Join(projectRoot(t), "bin")
	if err := os.RemoveAll(binDir); err != nil {
		t.Fatalf("failed to clean bin directory: %v", err)
	}
}
