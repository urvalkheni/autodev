package catalog

import (
	"testing"
)

func TestCatalogResolution(t *testing.T) {
	c, err := Load()
	if err != nil {
		t.Fatalf("failed to load catalog: %v", err)
	}

	// Test profile resolution
	resolved, err := c.ResolveProfile("web-dev")
	if err != nil {
		t.Fatalf("failed to resolve web-dev profile: %v", err)
	}

	if len(resolved) == 0 {
		t.Fatalf("resolved list is empty")
	}

	// Verify that dependencies come before dependants
	// e.g. nodejs must come before react
	nodeIdx := -1
	reactIdx := -1
	for i, pkg := range resolved {
		if pkg.ID == "nodejs" {
			nodeIdx = i
		} else if pkg.ID == "react" {
			reactIdx = i
		}
	}

	if nodeIdx == -1 {
		t.Errorf("nodejs not found in resolved profile")
	}
	if reactIdx == -1 {
		t.Errorf("react not found in resolved profile")
	}
	if nodeIdx > reactIdx {
		t.Errorf("nodejs resolved after react: nodejs at %d, react at %d", nodeIdx, reactIdx)
	}
}
