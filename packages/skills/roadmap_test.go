package skills

import (
	"testing"
)

func TestGetAvailableSkills(t *testing.T) {
	gen := New()
	skills := gen.GetAvailableSkills()

	if len(skills) == 0 {
		t.Fatal("Expected catalog to have skills, got 0")
	}

	// Verify it's sorted
	for i := 1; i < len(skills); i++ {
		if skills[i] < skills[i-1] {
			t.Errorf("Catalog keys not sorted: %s < %s", skills[i], skills[i-1])
		}
	}

	// Verify common skills are present
	foundGo := false
	foundNode := false
	for _, s := range skills {
		if s == "Go" {
			foundGo = true
		}
		if s == "Node.js" {
			foundNode = true
		}
	}

	if !foundGo {
		t.Error("Expected Go to be in available skills catalog")
	}
	if !foundNode {
		t.Error("Expected Node.js to be in available skills catalog")
	}
}

func TestGenerateRoadmap(t *testing.T) {
	gen := New()
	detected := []string{"Go", "Node.js"}
	roadmap := gen.Generate(detected)

	if len(roadmap.CurrentSkills) != 2 {
		t.Errorf("Expected 2 current skills, got %d", len(roadmap.CurrentSkills))
	}

	// Verify next steps and goals are populated
	if len(roadmap.NextSteps) == 0 {
		t.Error("Expected next steps to be recommended, got 0")
	}

	// Check fields are filled
	if roadmap.Title == "" {
		t.Error("Expected title to be non-empty")
	}
	if roadmap.GeneratedAt == "" {
		t.Error("Expected generated_at to be non-empty")
	}
}
