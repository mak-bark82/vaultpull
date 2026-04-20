package envnotify

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// Level represents the severity of a notification.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
)

// Event represents a single notification event.
type Event struct {
	Timestamp time.Time
	Level     Level
	Key       string
	Message   string
}

// Notifier collects and emits structured notification events.
type Notifier struct {
	writer io.Writer
	clock  func() time.Time
	events []Event
}

// New creates a Notifier that writes formatted events to w.
func New(w io.Writer) *Notifier {
	return &Notifier{
		writer: w,
		clock:  time.Now,
	}
}

// Notify records an event and writes it to the underlying writer.
func (n *Notifier) Notify(level Level, key, message string) {
	e := Event{
		Timestamp: n.clock(),
		Level:     level,
		Key:       key,
		Message:   message,
	}
	n.events = append(n.events, e)
	if n.writer != nil {
		fmt.Fprintf(n.writer, "[%s] %s key=%s msg=%s\n",
			e.Timestamp.Format(time.RFC3339), e.Level, e.Key, e.Message)
	}
}

// Events returns a copy of all recorded events.
func (n *Notifier) Events() []Event {
	out := make([]Event, len(n.events))
	copy(out, n.events)
	return out
}

// Summary returns a human-readable summary of events grouped by level.
func (n *Notifier) Summary() string {
	counts := map[Level]int{}
	for _, e := range n.events {
		counts[e.Level]++
	}
	parts := []string{}
	for _, lvl := range []Level{LevelInfo, LevelWarn, LevelError} {
		if c, ok := counts[lvl]; ok {
			parts = append(parts, fmt.Sprintf("%s=%d", lvl, c))
		}
	}
	if len(parts) == 0 {
		return "no events"
	}
	return strings.Join(parts, " ")
}
