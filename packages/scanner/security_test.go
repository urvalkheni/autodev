package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCleanVersion(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"^4.17.21", "4.17.21"},
		{"~1.2.3", "1.2.3"},
		{">=1.0.0 <2.0.0", "1.0.0"},
		{"v1.9.0", "1.9.0"},
		{"  1.5.0  ", "1.5.0"},
	}

	for _, test := range tests {
		result := CleanVersion(test.input)
		if result != test.expected {
			t.Errorf("CleanVersion(%q) = %q; expected %q", test.input, result, test.expected)
		}
	}
}

func TestParsePackageJSON(t *testing.T) {
	content := `{
		"dependencies": {
			"lodash": "^4.17.21"
		},
		"devDependencies": {
			"typescript": "5.0.4"
		}
	}`

	tmpDir, err := os.MkdirTemp("", "autodev-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "package.json")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	pkgs := parsePackageJSON(filePath)
	if len(pkgs) != 2 {
		t.Errorf("Expected 2 packages, got %d", len(pkgs))
	}

	foundLodash := false
	foundTS := false
	for _, p := range pkgs {
		if p.Name == "lodash" && p.Version == "4.17.21" && p.Ecosystem == "npm" {
			foundLodash = true
		}
		if p.Name == "typescript" && p.Version == "5.0.4" && p.Ecosystem == "npm" {
			foundTS = true
		}
	}

	if !foundLodash {
		t.Errorf("lodash package not parsed correctly: %+v", pkgs)
	}
	if !foundTS {
		t.Errorf("typescript package not parsed correctly: %+v", pkgs)
	}
}

func TestParseGoMod(t *testing.T) {
	content := `module testapp

go 1.22.2

require (
	github.com/spf13/cobra v1.8.0
	github.com/spf13/viper v1.18.2
)

require github.com/charmbracelet/lipgloss v0.9.1 // indirect
`

	tmpDir, err := os.MkdirTemp("", "autodev-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	pkgs := parseGoMod(filePath)
	if len(pkgs) != 3 {
		t.Errorf("Expected 3 packages, got %d: %+v", len(pkgs), pkgs)
	}

	foundCobra := false
	foundLipgloss := false
	for _, p := range pkgs {
		if p.Name == "github.com/spf13/cobra" && p.Version == "1.8.0" && p.Ecosystem == "Go" {
			foundCobra = true
		}
		if p.Name == "github.com/charmbracelet/lipgloss" && p.Version == "0.9.1" && p.Ecosystem == "Go" {
			foundLipgloss = true
		}
	}

	if !foundCobra {
		t.Errorf("cobra package not parsed correctly: %+v", pkgs)
	}
	if !foundLipgloss {
		t.Errorf("lipgloss package not parsed correctly: %+v", pkgs)
	}
}

func TestParseRequirements(t *testing.T) {
	content := `
requests==2.31.0
flask>=2.0.0
# some comment
numpy
`

	tmpDir, err := os.MkdirTemp("", "autodev-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "requirements.txt")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	pkgs := parseRequirements(filePath)
	if len(pkgs) != 2 {
		t.Errorf("Expected 2 packages, got %d: %+v", len(pkgs), pkgs)
	}

	foundRequests := false
	foundFlask := false
	for _, p := range pkgs {
		if p.Name == "requests" && p.Version == "2.31.0" && p.Ecosystem == "PyPI" {
			foundRequests = true
		}
		if p.Name == "flask" && p.Version == "2.0.0" && p.Ecosystem == "PyPI" {
			foundFlask = true
		}
	}

	if !foundRequests {
		t.Errorf("requests package not parsed correctly: %+v", pkgs)
	}
	if !foundFlask {
		t.Errorf("flask package not parsed correctly: %+v", pkgs)
	}
}
