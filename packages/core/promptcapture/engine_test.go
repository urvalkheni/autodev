package promptcapture

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestEngineLifecycle(t *testing.T) {
	// Setup temp directory to mimic project root
	tempDir, err := os.MkdirTemp("", "autodev-promptcapture-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a mock package.json to help root detection
	err = os.WriteFile(filepath.Join(tempDir, "package.json"), []byte(`{"name":"test-project"}`), 0644)
	if err != nil {
		t.Fatalf("failed to create package.json: %v", err)
	}

	engine, err := NewEngine(tempDir)
	if err != nil {
		t.Fatalf("failed to initialize Engine: %v", err)
	}

	if engine.Root != tempDir {
		t.Errorf("expected engine Root to be %s, got %s", tempDir, engine.Root)
	}

	// Verify dirs were created
	dirs := []string{"sessions", "prompts", "workflows", "analytics"}
	for _, d := range dirs {
		path := filepath.Join(tempDir, ".autodevs", d)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("directory %s was not created", path)
		}
	}

	// Start session
	err = engine.StartSession()
	if err != nil {
		t.Fatalf("failed to start session: %v", err)
	}

	if engine.Session == nil {
		t.Fatal("engine session should not be nil after StartSession")
	}

	sessionID := engine.Session.SessionID
	if !strings.Contains(sessionID, time.Now().Format("2006-01-02")) {
		t.Errorf("session ID should contain current date, got %s", sessionID)
	}

	// Track some simulated generated files and commands
	cmds := []ExecutedCommand{
		{
			Command:    "go",
			Args:       []string{"test", "./..."},
			ExitCode:   0,
			Stdout:     "PASS",
			DurationMs: 120,
			Timestamp:  time.Now(),
		},
	}
	files := []GeneratedFile{
		{
			FilePath:  "main.go",
			SizeBytes: 300,
			Action:    "created",
			Timestamp: time.Now(),
		},
	}

	err = engine.AddEvent("run tests and generate main", "tests passed and main.go created", cmds, files)
	if err != nil {
		t.Fatalf("failed to add prompt event: %v", err)
	}

	// Verify files were written
	jsonPath := filepath.Join(tempDir, ".autodevs", "sessions", sessionID+".json")
	if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
		t.Errorf("JSON session log not found at %s", jsonPath)
	}

	mdPath := filepath.Join(tempDir, ".autodevs", "sessions", sessionID+".md")
	if _, err := os.Stat(mdPath); os.IsNotExist(err) {
		t.Errorf("Markdown session log not found at %s", mdPath)
	}

	rootMDPath := filepath.Join(tempDir, "prompts.md")
	if _, err := os.Stat(rootMDPath); os.IsNotExist(err) {
		t.Errorf("Root prompts.md not found at %s", rootMDPath)
	}

	// Verify prompts.json log
	promptsLogPath := filepath.Join(tempDir, ".autodevs", "prompts", "prompts.json")
	if _, err := os.Stat(promptsLogPath); os.IsNotExist(err) {
		t.Errorf("Prompts index prompts.json not found")
	}

	// End Session
	err = engine.EndSession()
	if err != nil {
		t.Fatalf("failed to end session: %v", err)
	}

	if engine.Session != nil {
		t.Errorf("session should be nil after EndSession")
	}
}

func TestOfflineQueue(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "autodev-offline-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	engine, err := NewEngine(tempDir)
	if err != nil {
		t.Fatalf("failed to initialize Engine: %v", err)
	}

	// Create test event
	ev := Event{
		Type:      EventPromptCaptured,
		Timestamp: time.Now(),
		Payload: PromptEvent{
			ID:     "test_id",
			Prompt: "hello",
		},
	}

	err = engine.QueueOfflineEvent(ev)
	if err != nil {
		t.Fatalf("failed to queue offline event: %v", err)
	}

	// Verify queue.json has the event
	queuePath := filepath.Join(tempDir, ".autodevs", "analytics", "queue.json")
	data, err := os.ReadFile(queuePath)
	if err != nil {
		t.Fatalf("failed to read queue.json: %v", err)
	}

	var queue []Event
	if err := json.Unmarshal(data, &queue); err != nil {
		t.Fatalf("failed to unmarshal queue.json: %v", err)
	}

	if len(queue) != 1 {
		t.Fatalf("expected queue length of 1, got %d", len(queue))
	}

	if queue[0].Type != EventPromptCaptured {
		t.Errorf("expected event type %s, got %s", EventPromptCaptured, queue[0].Type)
	}
}
