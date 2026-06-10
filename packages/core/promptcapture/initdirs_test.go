package promptcapture

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitDirsRalphLoopScaffold(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "autodev-initdirs-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	if err := InitDirs(tempDir); err != nil {
		t.Fatalf("InitDirs failed: %v", err)
	}

	ralphDirs := []string{
		"todo", "completed", "failed", "logs", "plans", "loops", "memory", "context", "agents",
	}
	for _, d := range ralphDirs {
		path := filepath.Join(tempDir, ".autodevs", d)
		info, err := os.Stat(path)
		if os.IsNotExist(err) {
			t.Errorf("Ralph Loop directory missing: %s", path)
			continue
		}
		if err != nil {
			t.Errorf("stat %s: %v", path, err)
			continue
		}
		if !info.IsDir() {
			t.Errorf("expected %s to be a directory", path)
		}
	}
}
