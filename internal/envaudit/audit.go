package envaudit

import (
	"fmt"
	"strings"
	"time"
)

// EventKind describes the type of audit event.
type EventKind string

const (
	EventAdded   EventKind = "added"
	EventChanged EventKind = "changed"
	EventRemoved EventKind = "removed"
	EventSkipped EventKind = "skipped"
)

// Event represents a single auditable change to an env key.
type Event struct {
	Timestamp time.Time
	Kind      EventKind
	Key       string
	EnvFile   string
	Note      string
}

// Recorder collects audit events during a sync run.
type Recorder struct {
	events []Event
	clock  func() time.Time
}

// New returns a new Recorder. If clock is nil, time.Now is used.
func New(clock func() time.Time) *Recorder {
	if clock == nil {
		clock = time.Now
	}
	return &Recorder{clock: clock}
}

// Record appends an audit event.
func (r *Recorder) Record(kind EventKind, key, envFile, note string) {
	r.events = append(r.events, Event{
		Timestamp: r.clock(),
		Kind:      kind,
		Key:       key,
		EnvFile:   envFile,
		Note:      note,
	})
}

// Events returns all recorded events.
func (r *Recorder) Events() []Event {
	out := make([]Event, len(r.events))
	copy(out, r.events)
	return out
}

// Summary returns a human-readable summary of all events.
func (r *Recorder) Summary() string {
	if len(r.events) == 0 {
		return "no audit events recorded"
	}
	var sb strings.Builder
	for _, e := range r.events {
		line := fmt.Sprintf("[%s] %s %s (%s)",
			e.Timestamp.Format(time.RFC3339), e.Kind, e.Key, e.EnvFile)
		if e.Note != "" {
			line += " — " + e.Note
		}
		sb.WriteString(line + "\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}
