package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScan_NodeProject(t *testing.T) {
	// Create temp directory simulating a Next.js project
	dir := t.TempDir()

	files := map[string]string{
		"package.json":   `{"name":"test","dependencies":{"react":"^18","next":"^15"}}`,
		"pnpm-lock.yaml": "",
		"next.config.js": "module.exports = {}",
		"Dockerfile":     "FROM node:20",
	}
	for name, content := range files {
		path := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	s := New(dir)
	result, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan() error: %v", err)
	}

	assertContains(t, result.Languages, "Node.js")
	assertContains(t, result.Frameworks, "Next.js")
	assertContains(t, result.Frameworks, "React")
	assertContains(t, result.PackageManagers, "pnpm")

	if !result.HasDocker {
		t.Error("expected HasDocker=true")
	}
}

func TestScan_GoProject(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/app\ngo 1.22\n"), 0644); err != nil {
		t.Fatal(err)
	}

	s := New(dir)
	result, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan() error: %v", err)
	}

	assertContains(t, result.Languages, "Go")
}

func TestScan_PythonProject(t *testing.T) {
	dir := t.TempDir()
	for _, f := range []string{"requirements.txt", "poetry.lock"} {
		if err := os.WriteFile(filepath.Join(dir, f), []byte(""), 0644); err != nil {
			t.Fatal(err)
		}
	}

	s := New(dir)
	result, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan() error: %v", err)
	}

	assertContains(t, result.Languages, "Python")
	assertContains(t, result.PackageManagers, "Poetry")
}

func assertContains(t *testing.T, slice []string, val string) {
	t.Helper()
	for _, s := range slice {
		if s == val {
			return
		}
	}
	t.Errorf("expected %q to contain %q, got %v", "slice", val, slice)
}
