package envwatch

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// FileState holds the last known checksum and modification time of a file.
type FileState struct {
	Path     string
	Checksum string
	ModTime  time.Time
}

// ChangeEvent is emitted when a watched file changes.
type ChangeEvent struct {
	Path    string
	OldHash string
	NewHash string
}

// Watcher polls a set of env files and emits ChangeEvents when content changes.
type Watcher struct {
	interval time.Duration
	mu       sync.Mutex
	states   map[string]FileState
	Events   chan ChangeEvent
	stop     chan struct{}
}

// New creates a Watcher that polls files at the given interval.
func New(interval time.Duration) *Watcher {
	return &Watcher{
		interval: interval,
		states:   make(map[string]FileState),
		Events:   make(chan ChangeEvent, 16),
		stop:     make(chan struct{}),
	}
}

// Watch registers a file path to be monitored.
func (w *Watcher) Watch(path string) error {
	state, err := snapshot(path)
	if err != nil {
		return fmt.Errorf("envwatch: initial snapshot of %q: %w", path, err)
	}
	w.mu.Lock()
	w.states[path] = state
	w.mu.Unlock()
	return nil
}

// Start begins polling in the background. Call Stop to terminate.
func (w *Watcher) Start() {
	ticker := time.NewTicker(w.interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				w.poll()
			case <-w.stop:
				return
			}
		}
	}()
}

// Stop halts the polling goroutine.
func (w *Watcher) Stop() {
	close(w.stop)
}

func (w *Watcher) poll() {
	w.mu.Lock()
	defer w.mu.Unlock()
	for path, old := range w.states {
		current, err := snapshot(path)
		if err != nil {
			continue
		}
		if current.Checksum != old.Checksum {
			w.Events <- ChangeEvent{Path: path, OldHash: old.Checksum, NewHash: current.Checksum}
			w.states[path] = current
		}
	}
}

func snapshot(path string) (FileState, error) {
	f, err := os.Open(path)
	if err != nil {
		return FileState{}, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return FileState{}, err
	}
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return FileState{}, err
	}
	return FileState{
		Path:     path,
		Checksum: hex.EncodeToString(h.Sum(nil)),
		ModTime:  info.ModTime(),
	}, nil
}
