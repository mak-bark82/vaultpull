package envnotify

import "github.com/dstotijn/vaultpull/internal/diff"

// Dispatcher maps diff.ChangeType values to notification levels and
// emits events for each change in a diff result set.
type Dispatcher struct {
	notifier *Notifier
}

// NewDispatcher wraps a Notifier so that diff results can be dispatched
// as structured notification events.
func NewDispatcher(n *Notifier) *Dispatcher {
	return &Dispatcher{notifier: n}
}

// Dispatch converts a slice of diff.Change values into Notifier events.
func (d *Dispatcher) Dispatch(changes []diff.Change) {
	for _, c := range changes {
		level, msg := levelForChange(c)
		d.notifier.Notify(level, c.Key, msg)
	}
}

func levelForChange(c diff.Change) (Level, string) {
	switch c.Type {
	case diff.Added:
		return LevelInfo, "key added"
	case diff.Removed:
		return LevelWarn, "key removed"
	case diff.Changed:
		return LevelWarn, "value changed"
	default:
		return LevelInfo, "unchanged"
	}
}
