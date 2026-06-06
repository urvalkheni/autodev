package promptcapture

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Engine struct {
	Root            string
	Session         *SessionLog
	Emitter         *EventEmitter
	DevMentorAPIURL string
}

func FindProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir
		}
		if _, err := os.Stat(filepath.Join(dir, "package.json")); err == nil {
			return dir
		}
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	wd, _ := os.Getwd()
	return wd
}

func InitDirs(root string) error {
	dirs := []string{
		filepath.Join(root, ".autodevs"),
		filepath.Join(root, ".autodevs", "sessions"),
		filepath.Join(root, ".autodevs", "prompts"),
		filepath.Join(root, ".autodevs", "workflows"),
		filepath.Join(root, ".autodevs", "analytics"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return err
		}
	}
	return nil
}

func GetNextSessionID(root string) string {
	dateStr := time.Now().Format("2006-01-02")
	sessionsDir := filepath.Join(root, ".autodevs", "sessions")
	idx := 1
	for {
		name := fmt.Sprintf("%s-%d.md", dateStr, idx)
		path := filepath.Join(sessionsDir, name)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return fmt.Sprintf("%s-%d", dateStr, idx)
		}
		idx++
	}
}

func DetectProjectMetadata(root string) ProjectMetadata {
	meta := ProjectMetadata{
		ProjectName: filepath.Base(root),
		Path:        root,
		Branch:      "unknown",
		Commit:      "unknown",
		Languages:   []string{},
		Frameworks:  []string{},
	}

	// 1. Git branch
	cmdBranch := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmdBranch.Dir = root
	if out, err := cmdBranch.Output(); err == nil {
		meta.Branch = strings.TrimSpace(string(out))
	}

	// 2. Git commit
	cmdCommit := exec.Command("git", "rev-parse", "--short", "HEAD")
	cmdCommit.Dir = root
	if out, err := cmdCommit.Output(); err == nil {
		meta.Commit = strings.TrimSpace(string(out))
	}

	// 3. Scan root files for languages & frameworks
	files, err := os.ReadDir(root)
	if err == nil {
		hasJS := false
		hasTS := false
		hasGo := false
		hasPython := false
		hasRust := false
		hasNext := false
		hasReact := false

		for _, f := range files {
			name := f.Name()
			if name == "go.mod" {
				hasGo = true
			} else if name == "package.json" {
				hasJS = true
				// Check inside package.json for react/next
				pkgData, err := os.ReadFile(filepath.Join(root, name))
				if err == nil {
					pkgStr := string(pkgData)
					if strings.Contains(pkgStr, `"next"`) {
						hasNext = true
					}
					if strings.Contains(pkgStr, `"react"`) {
						hasReact = true
					}
					if strings.Contains(pkgStr, `"typescript"`) {
						hasTS = true
					}
				}
			} else if name == "Cargo.toml" {
				hasRust = true
			} else if name == "requirements.txt" || name == "Pipfile" || name == "pyproject.toml" {
				hasPython = true
			} else if name == "tsconfig.json" {
				hasTS = true
			} else if name == "next.config.js" || name == "next.config.mjs" {
				hasNext = true
			}
		}

		if hasGo {
			meta.Languages = append(meta.Languages, "Go")
		}
		if hasRust {
			meta.Languages = append(meta.Languages, "Rust")
		}
		if hasPython {
			meta.Languages = append(meta.Languages, "Python")
		}
		if hasTS {
			meta.Languages = append(meta.Languages, "TypeScript")
		} else if hasJS {
			meta.Languages = append(meta.Languages, "JavaScript")
		}

		if hasNext {
			meta.Frameworks = append(meta.Frameworks, "Next.js")
		}
		if hasReact && !hasNext {
			meta.Frameworks = append(meta.Frameworks, "React")
		}
	}

	return meta
}

func NewEngine(root string) (*Engine, error) {
	if root == "" {
		root = FindProjectRoot()
	}
	if err := InitDirs(root); err != nil {
		return nil, err
	}
	emitter := NewEventEmitter()

	apiURL := os.Getenv("DEVMENTOR_API_URL")
	if apiURL == "" {
		apiURL = "https://api.devmentor.co/v1/events"
	}

	engine := &Engine{
		Root:            root,
		Emitter:         emitter,
		DevMentorAPIURL: apiURL,
	}

	engine.setupDefaultListeners()

	return engine, nil
}

func (e *Engine) setupDefaultListeners() {
	e.Emitter.On(EventPromptCaptured, func(event Event) {
		err := e.sendToDevMentor(event)
		if err != nil {
			_ = e.QueueOfflineEvent(event)
		}
	})
	e.Emitter.On(EventSessionStarted, func(event Event) {
		err := e.sendToDevMentor(event)
		if err != nil {
			_ = e.QueueOfflineEvent(event)
		}
	})
	e.Emitter.On(EventSessionEnded, func(event Event) {
		err := e.sendToDevMentor(event)
		if err != nil {
			_ = e.QueueOfflineEvent(event)
		}
	})
}

func (e *Engine) StartSession() error {
	id := GetNextSessionID(e.Root)
	meta := DetectProjectMetadata(e.Root)

	e.Session = &SessionLog{
		SessionID: id,
		StartTime: time.Now(),
		Metadata:  meta,
		Events:    []PromptEvent{},
	}

	e.Emitter.Emit(Event{
		Type:      EventSessionStarted,
		Payload:   e.Session,
		Timestamp: e.Session.StartTime,
	})

	return e.SaveSession()
}

func (e *Engine) LoadLatestSession() error {
	sessionsDir := filepath.Join(e.Root, ".autodevs", "sessions")
	files, err := os.ReadDir(sessionsDir)
	if err != nil {
		return err
	}

	var latestFile string
	var latestTime time.Time

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".json") {
			info, err := f.Info()
			if err == nil && info.ModTime().After(latestTime) {
				latestTime = info.ModTime()
				latestFile = f.Name()
			}
		}
	}

	if latestFile == "" {
		return fmt.Errorf("no active session found, start one with 'autodev chat'")
	}

	data, err := os.ReadFile(filepath.Join(sessionsDir, latestFile))
	if err != nil {
		return err
	}

	var log SessionLog
	if err := json.Unmarshal(data, &log); err != nil {
		return err
	}

	e.Session = &log
	return nil
}

func (e *Engine) AddEvent(prompt string, response string, cmds []ExecutedCommand, files []GeneratedFile) error {
	if e.Session == nil {
		if err := e.StartSession(); err != nil {
			return err
		}
	}

	eventID := fmt.Sprintf("prompt_%d", len(e.Session.Events)+1)
	ev := PromptEvent{
		ID:               eventID,
		Timestamp:        time.Now(),
		Prompt:           prompt,
		Response:         response,
		ExecutedCommands: cmds,
		GeneratedFiles:   files,
		Metadata:         e.Session.Metadata,
	}

	e.Session.Events = append(e.Session.Events, ev)

	// Append to prompts log
	_ = e.AppendToPromptsLog(prompt)

	e.Emitter.Emit(Event{
		Type:      EventPromptCaptured,
		Payload:   ev,
		Timestamp: ev.Timestamp,
	})

	return e.SaveSession()
}

func (e *Engine) AppendToPromptsLog(prompt string) error {
	logPath := filepath.Join(e.Root, ".autodevs", "prompts", "prompts.json")
	var prompts []GlobalPrompt

	if _, err := os.Stat(logPath); err == nil {
		data, err := os.ReadFile(logPath)
		if err == nil {
			_ = json.Unmarshal(data, &prompts)
		}
	}

	prompts = append(prompts, GlobalPrompt{
		Timestamp: time.Now(),
		SessionID: e.Session.SessionID,
		Prompt:    prompt,
	})

	data, err := json.MarshalIndent(prompts, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(logPath, data, 0644)
}

func (e *Engine) SaveSession() error {
	if e.Session == nil {
		return fmt.Errorf("no active session")
	}

	sessionsDir := filepath.Join(e.Root, ".autodevs", "sessions")

	// Save JSON version for replaying/indexing
	jsonPath := filepath.Join(sessionsDir, fmt.Sprintf("%s.json", e.Session.SessionID))
	jsonData, err := json.MarshalIndent(e.Session, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(jsonPath, jsonData, 0644); err != nil {
		return err
	}

	// Generate Markdown
	mdContent := e.Session.ToMarkdown()

	// Save session Markdown
	mdPath := filepath.Join(sessionsDir, fmt.Sprintf("%s.md", e.Session.SessionID))
	if err := os.WriteFile(mdPath, []byte(mdContent), 0644); err != nil {
		return err
	}

	// Update root-level prompts.md
	rootMDPath := filepath.Join(e.Root, "prompts.md")
	if err := os.WriteFile(rootMDPath, []byte(mdContent), 0644); err != nil {
		return err
	}

	return nil
}

func (e *Engine) EndSession() error {
	if e.Session == nil {
		return nil
	}

	e.Session.EndTime = time.Now()
	e.Emitter.Emit(Event{
		Type:      EventSessionEnded,
		Payload:   e.Session,
		Timestamp: e.Session.EndTime,
	})

	err := e.SaveSession()
	e.Session = nil
	return err
}

func (s *SessionLog) ToMarkdown() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Autodevs Session - %s\n\n", s.StartTime.Format("2006-01-02 15:04:05")))
	sb.WriteString("## Project Metadata\n")
	sb.WriteString(fmt.Sprintf("- **Project Name:** %s\n", s.Metadata.ProjectName))
	sb.WriteString(fmt.Sprintf("- **Path:** `%s`\n", s.Metadata.Path))
	sb.WriteString(fmt.Sprintf("- **Git Branch:** `%s`\n", s.Metadata.Branch))
	sb.WriteString(fmt.Sprintf("- **Last Commit:** `%s`\n", s.Metadata.Commit))
	if len(s.Metadata.Languages) > 0 {
		sb.WriteString(fmt.Sprintf("- **Languages:** %s\n", strings.Join(s.Metadata.Languages, ", ")))
	}
	if len(s.Metadata.Frameworks) > 0 {
		sb.WriteString(fmt.Sprintf("- **Frameworks:** %s\n", strings.Join(s.Metadata.Frameworks, ", ")))
	}
	sb.WriteString("\n---\n\n")

	for i, ev := range s.Events {
		sb.WriteString(fmt.Sprintf("## Prompt %d (Captured: %s)\n", i+1, ev.Timestamp.Format("15:04:05")))
		sb.WriteString("### User\n")
		sb.WriteString(ev.Prompt)
		sb.WriteString("\n\n")
		sb.WriteString("### AI\n")
		sb.WriteString(ev.Response)
		sb.WriteString("\n\n")

		if len(ev.GeneratedFiles) > 0 {
			sb.WriteString("#### Generated Files\n")
			for _, f := range ev.GeneratedFiles {
				sb.WriteString(fmt.Sprintf("- `%s` (%d bytes) - %s\n", f.FilePath, f.SizeBytes, f.Action))
			}
			sb.WriteString("\n")
		}

		if len(ev.ExecutedCommands) > 0 {
			sb.WriteString("#### Executed Commands\n")
			for _, cmd := range ev.ExecutedCommands {
				sb.WriteString(fmt.Sprintf("- `%s %s` (Exit Code: %d, Duration: %.2fs)\n", 
					cmd.Command, strings.Join(cmd.Args, " "), cmd.ExitCode, float64(cmd.DurationMs)/1000.0))
			}
			sb.WriteString("\n")
		}

		sb.WriteString("---\n\n")
	}

	return sb.String()
}

func (e *Engine) QueueOfflineEvent(event Event) error {
	queuePath := filepath.Join(e.Root, ".autodevs", "analytics", "queue.json")
	var queue []Event

	if _, err := os.Stat(queuePath); err == nil {
		data, err := os.ReadFile(queuePath)
		if err == nil {
			_ = json.Unmarshal(data, &queue)
		}
	}

	queue = append(queue, event)

	data, err := json.MarshalIndent(queue, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(queuePath, data, 0644)
}

func (e *Engine) SyncOfflineEvents() (int, error) {
	queuePath := filepath.Join(e.Root, ".autodevs", "analytics", "queue.json")
	if _, err := os.Stat(queuePath); os.IsNotExist(err) {
		return 0, nil
	}

	data, err := os.ReadFile(queuePath)
	if err != nil {
		return 0, err
	}

	var queue []Event
	if err := json.Unmarshal(data, &queue); err != nil {
		return 0, err
	}

	if len(queue) == 0 {
		return 0, nil
	}

	var remaining []Event
	successCount := 0

	for _, ev := range queue {
		err := e.sendToDevMentor(ev)
		if err != nil {
			remaining = append(remaining, ev)
		} else {
			successCount++
		}
	}

	if len(remaining) == 0 {
		_ = os.Remove(queuePath)
	} else {
		newData, err := json.MarshalIndent(remaining, "", "  ")
		if err == nil {
			_ = os.WriteFile(queuePath, newData, 0644)
		}
	}

	return successCount, nil
}

func (e *Engine) sendToDevMentor(event Event) error {
	payload := DevMentorEventPayload{
		Event:     string(event.Type),
		Timestamp: event.Timestamp.Format(time.RFC3339),
		Data:      event.Payload,
	}
	
	if promptEv, ok := event.Payload.(PromptEvent); ok {
		payload.SessionID = promptEv.ID
	} else if sessionEv, ok := event.Payload.(*SessionLog); ok {
		payload.SessionID = sessionEv.SessionID
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", e.DevMentorAPIURL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	if apiKey := os.Getenv("DEVMENTOR_API_KEY"); apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("DevMentor API returned status: %d", resp.StatusCode)
	}

	return nil
}
