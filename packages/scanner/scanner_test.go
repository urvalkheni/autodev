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

func TestScan_Monorepo(t *testing.T) {
	dir := t.TempDir()

	// Setup a monorepo workspace with 2 subprojects: cli (Go) and web (Next.js)
	files := map[string]string{
		"package.json":                `{"name":"monorepo"}`,
		"packages/cli/go.mod":         "module example.com/cli\ngo 1.22\n",
		"apps/web/package.json":       `{"name":"web","dependencies":{"react":"^18","next":"^15"}}`,
		"apps/web/next.config.js":     "module.exports = {}",
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

	// Verify we detected Go and Node/Next
	assertContains(t, result.Languages, "Go")
	assertContains(t, result.Languages, "Node.js")
	assertContains(t, result.Frameworks, "Next.js")

	// Verify subprojects are mapped
	if len(result.Projects) != 2 {
		t.Fatalf("expected 2 subprojects, got %d: %+v", len(result.Projects), result.Projects)
	}

	var foundCli, foundWeb bool
	for _, proj := range result.Projects {
		if proj.Path == "packages/cli" {
			foundCli = true
			var hasGo bool
			for _, tech := range proj.Technologies {
				if tech.Name == "Go" {
					hasGo = true
				}
			}
			if !hasGo {
				t.Errorf("cli subproject missing Go technology: %+v", proj)
			}
		}
		if proj.Path == "apps/web" {
			foundWeb = true
			var hasNext bool
			for _, tech := range proj.Technologies {
				if tech.Name == "Next.js" {
					hasNext = true
				}
			}
			if !hasNext {
				t.Errorf("web subproject missing Next.js technology: %+v", proj)
			}
		}
	}

	if !foundCli {
		t.Error("expected to find cli subproject in packages/cli")
	}
	if !foundWeb {
		t.Error("expected to find web subproject in apps/web")
	}
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
