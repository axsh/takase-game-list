package integration

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestBinaryExists_Integration(t *testing.T) {
	binPath := filepath.Join(projectRoot(t), "bin", "game-list.exe")
	info, err := os.Stat(binPath)
	if err != nil {
		t.Fatalf("binary not found at %s (run build script first): %v", binPath, err)
	}
	if info.Size() == 0 {
		t.Fatalf("binary at %s is empty", binPath)
	}
}
func projectRoot(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to resolve caller path")
	}
	return filepath.Dir(filepath.Dir(filepath.Dir(filename)))
}
