package promptcapture

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIsSensitive(t *testing.T) {
	tests := []struct {
		prompt string
		want   bool
	}{
		{"create a beautiful layout button", false},
		{"set mysql_pwd=mysecurepassword", true},
		{"find code containing private_key", true},
		{"generate random password string", true},
		{"optimize postgresql query", false},
	}

	for _, tt := range tests {
		got := IsSensitive(tt.prompt)
		if got != tt.want {
			t.Errorf("IsSensitive(%q) = %v; want %v", tt.prompt, got, tt.want)
		}
	}
}

func TestAppendToPromptsMD(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "autodev-capture-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	prompt := "optimize the fibonacci function"
	err = AppendToPromptsMD(tempDir, prompt)
	if err != nil {
		t.Fatalf("failed to append prompt: %v", err)
	}

	// Read prompts.md and verify contents
	mdPath := filepath.Join(tempDir, ".autodevs", "prompts.md")
	data, err := os.ReadFile(mdPath)
	if err != nil {
		t.Fatalf("failed to read prompts.md: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "# Prompt History") {
		t.Errorf("expected header '# Prompt History', got:\n%s", content)
	}
	if !strings.Contains(content, prompt) {
		t.Errorf("expected prompt text %q, got:\n%s", prompt, content)
	}

	// Test duplicate insertion
	err = AppendToPromptsMD(tempDir, prompt)
	if err != nil {
		t.Fatalf("failed second append: %v", err)
	}

	data2, err := os.ReadFile(mdPath)
	if err != nil {
		t.Fatalf("failed to read prompts.md: %v", err)
	}

	// Counting occurrences of the prompt in the file
	count := strings.Count(string(data2), prompt)
	if count != 1 {
		t.Errorf("expected prompt to be deduplicated (count 1), got %d occurrences", count)
	}
}

func TestExtractPromptFromCmdline(t *testing.T) {
	tests := []struct {
		cmdline []string
		want    string
	}{
		{[]string{"gemini", "create a react component"}, "create a react component"},
		{[]string{"claude", "--verbose", "optimize db queries"}, "optimize db queries"},
		{[]string{"copilot", "-v"}, ""},
		{[]string{"agy", "fix the broken imports"}, "fix the broken imports"},
	}

	for _, tt := range tests {
		got := ExtractPromptFromCmdline(tt.cmdline)
		if got != tt.want {
			t.Errorf("ExtractPromptFromCmdline(%v) = %q; want %q", tt.cmdline, got, tt.want)
		}
	}
}

func TestStdinProxy(t *testing.T) {
	var captured []string
	onLine := func(line string) {
		captured = append(captured, line)
	}

	input := "hello world\nignoreme\nthis is a captured prompt\n"
	src := strings.NewReader(input)
	var dst bytes.Buffer

	proxy := NewStdinProxy(src, onLine)

	_, err := io.Copy(&dst, proxy)
	if err != nil {
		t.Fatalf("failed to copy: %v", err)
	}

	// Output destination should match source input
	if dst.String() != input {
		t.Errorf("expected stdout copy %q, got %q", input, dst.String())
	}

	// "ignoreme" is a single word, should be skipped. "hello world" and "this is a captured prompt" should be logged.
	if len(captured) != 2 {
		t.Fatalf("expected 2 captured lines, got %d: %v", len(captured), captured)
	}

	if captured[0] != "hello world" {
		t.Errorf("expected first line 'hello world', got %q", captured[0])
	}
	if captured[1] != "this is a captured prompt" {
		t.Errorf("expected second line 'this is a captured prompt', got %q", captured[1])
	}
}

func TestImportAntigravityPrompts(t *testing.T) {
	// Create a temp directory to act as the user's HOME directory
	mockHome, err := os.MkdirTemp("", "mock-home-")
	if err != nil {
		t.Fatalf("failed to create mock home: %v", err)
	}
	defer os.RemoveAll(mockHome)

	// Save original HOME env var and restore it at the end of the test
	origHome := os.Getenv("HOME")
	defer os.Setenv("HOME", origHome)

	// Set HOME to mockHome
	if err := os.Setenv("HOME", mockHome); err != nil {
		t.Fatalf("failed to set HOME env: %v", err)
	}

	// Create the folder structure: <mockHome>/.gemini/antigravity/brain/session1/.system_generated/logs/
	sessionDir := filepath.Join(mockHome, ".gemini", "antigravity", "brain", "session1")
	logsDir := filepath.Join(sessionDir, ".system_generated", "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		t.Fatalf("failed to create logs dir: %v", err)
	}

	// Write overview.txt containing a mix of USER_INPUT from USER_EXPLICIT, USER_INPUT from MODEL, and MODEL response
	overviewContent := `{"source":"USER_EXPLICIT","type":"USER_INPUT","content":"<USER_REQUEST>\nImplement a prompt capture engine.\n</USER_REQUEST>"}
{"source":"MODEL","type":"PLANNER_RESPONSE","content":"Implementing the engine."}
{"source":"MODEL","type":"USER_INPUT","content":"This is a model input that should be ignored."}
{"source":"USER_EXPLICIT","type":"USER_INPUT","content":"<USER_REQUEST>\nShow prompts history.\n</USER_REQUEST>"}
`
	overviewPath := filepath.Join(logsDir, "overview.txt")
	if err := os.WriteFile(overviewPath, []byte(overviewContent), 0644); err != nil {
		t.Fatalf("failed to write overview.txt: %v", err)
	}

	// Create a temp workspaceRoot directory
	workspaceRoot, err := os.MkdirTemp("", "mock-workspace-")
	if err != nil {
		t.Fatalf("failed to create mock workspace: %v", err)
	}
	defer os.RemoveAll(workspaceRoot)

	// Call ImportAntigravityPrompts
	err = ImportAntigravityPrompts(workspaceRoot)
	if err != nil {
		t.Fatalf("ImportAntigravityPrompts failed: %v", err)
	}

	// Read the written prompts.md in the workspaceRoot/.autodevs/prompts.md
	promptsPath := filepath.Join(workspaceRoot, ".autodevs", "prompts.md")
	data, err := os.ReadFile(promptsPath)
	if err != nil {
		t.Fatalf("failed to read prompts.md: %v", err)
	}

	content := string(data)

	// Assertions:
	// - Should contain "Implement a prompt capture engine."
	// - Should contain "Show prompts history."
	// - Should NOT contain "Implementing the engine."
	// - Should NOT contain "This is a model input that should be ignored."
	if !strings.Contains(content, "Implement a prompt capture engine.") {
		t.Errorf("expected user input prompt to be imported, but was not. Got:\n%s", content)
	}
	if !strings.Contains(content, "Show prompts history.") {
		t.Errorf("expected second user input prompt to be imported, but was not. Got:\n%s", content)
	}
	if strings.Contains(content, "Implementing the engine.") {
		t.Errorf("expected model response to be filtered out, but was found. Got:\n%s", content)
	}
	if strings.Contains(content, "This is a model input that should be ignored.") {
		t.Errorf("expected model source USER_INPUT to be filtered out, but was found. Got:\n%s", content)
	}
}

func TestCleanPrompt(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			"<USER_REQUEST>\nHello World\n</USER_REQUEST>",
			"Hello World",
		},
		{
			"<USER_REQUEST>\nHello World\n<truncated 123 bytes>",
			"Hello World\n<truncated 123 bytes>",
		},
		{
			"Some raw prompt without tags",
			"Some raw prompt without tags",
		},
		{
			"<USER_REQUEST>\nImplement this\n</USER_REQUEST>\n<USER_INFORMATION>\nSome user system info\n</USER_INFORMATION>",
			"Implement this",
		},
		{
			"<USER_REQUEST>\nImplement that\n<USER_INFORMATION>\nSome user system info",
			"Implement that",
		},
	}

	for _, tt := range tests {
		got := cleanPrompt(tt.input)
		if got != tt.want {
			t.Errorf("cleanPrompt(%q) = %q; want %q", tt.input, got, tt.want)
		}
	}
}

func TestIsDuplicate(t *testing.T) {
	content := `# Prompt History

## 2026-06-05 17:05:40

hello world

---

## 2026-06-05 17:10:00

first line
second line

---
`

	tests := []struct {
		prompt string
		want   bool
	}{
		{"hello world", true},
		{"hello world\n", true},
		{"first line\nsecond line", true},
		{"first line\r\nsecond line", true},
		{"different prompt", false},
		{"hello", false},
	}

	for _, tt := range tests {
		got := isDuplicate(content, tt.prompt)
		if got != tt.want {
			t.Errorf("isDuplicate(content, %q) = %v; want %v", tt.prompt, got, tt.want)
		}
	}
}
