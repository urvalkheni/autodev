package promptcapture

import (
	"sync"
	"time"
)

type EventType string

const (
	EventPromptCaptured EventType = "prompt.captured"
	EventCommandExecuted EventType = "command.executed"
	EventFileGenerated   EventType = "file.generated"
	EventSessionStarted  EventType = "session.started"
	EventSessionEnded    EventType = "session.ended"
)

type Event struct {
	Type      EventType   `json:"type"`
	Payload   interface{} `json:"payload"`
	Timestamp time.Time   `json:"timestamp"`
}

type ProjectMetadata struct {
	ProjectName string   `json:"project_name"`
	Path        string   `json:"path"`
	Branch      string   `json:"branch"`
	Commit      string   `json:"commit"`
	Languages   []string `json:"languages"`
	Frameworks  []string `json:"frameworks"`
}

type ExecutedCommand struct {
	Command    string    `json:"command"`
	Args       []string  `json:"args"`
	ExitCode   int       `json:"exit_code"`
	Stdout     string    `json:"stdout"`
	Stderr     string    `json:"stderr"`
	DurationMs int64     `json:"duration_ms"`
	Timestamp  time.Time `json:"timestamp"`
}

type GeneratedFile struct {
	FilePath  string    `json:"file_path"`
	SizeBytes int64     `json:"size_bytes"`
	Action    string    `json:"action"` // "created" or "modified"
	Timestamp time.Time `json:"timestamp"`
}

type PromptEvent struct {
	ID               string            `json:"id"`
	Timestamp        time.Time         `json:"timestamp"`
	Prompt           string            `json:"prompt"`
	Response         string            `json:"response"`
	ExecutedCommands []ExecutedCommand `json:"executed_commands,omitempty"`
	GeneratedFiles   []GeneratedFile   `json:"generated_files,omitempty"`
	Metadata         ProjectMetadata   `json:"metadata"`
}

type SessionLog struct {
	SessionID string          `json:"session_id"`
	StartTime time.Time       `json:"start_time"`
	EndTime   time.Time       `json:"end_time"`
	Events    []PromptEvent   `json:"events"`
	Metadata  ProjectMetadata `json:"metadata"`
}

type GlobalPrompt struct {
	Timestamp time.Time `json:"timestamp"`
	SessionID string    `json:"session_id"`
	Prompt    string    `json:"prompt"`
}

type DevMentorEventPayload struct {
	Event     string      `json:"event"`
	SessionID string      `json:"session_id"`
	Timestamp string      `json:"timestamp"`
	Data      interface{} `json:"data"`
}

type EventEmitter struct {
	listeners map[EventType][]func(Event)
	mu        sync.RWMutex
}

func NewEventEmitter() *EventEmitter {
	return &EventEmitter{
		listeners: make(map[EventType][]func(Event)),
	}
}

func (ee *EventEmitter) On(eventType EventType, listener func(Event)) {
	ee.mu.Lock()
	defer ee.mu.Unlock()
	ee.listeners[eventType] = append(ee.listeners[eventType], listener)
}

func (ee *EventEmitter) Emit(event Event) {
	ee.mu.RLock()
	defer ee.mu.RUnlock()
	if list, ok := ee.listeners[event.Type]; ok {
		for _, listener := range list {
			go listener(event)
		}
	}
}
