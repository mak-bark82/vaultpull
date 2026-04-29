package envlease

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

// Lease represents a time-bounded claim on a secret key.
type Lease struct {
	Key       string    `json:"key"`
	ExpiresAt time.Time `json:"expires_at"`
	Owner     string    `json:"owner"`
}

// IsExpired returns true if the lease has passed its expiry time.
func (l Lease) IsExpired(now time.Time) bool {
	return now.After(l.ExpiresAt)
}

// Manager tracks active leases for env keys.
type Manager struct {
	mu     sync.RWMutex
	leases map[string]Lease
	path   string
	clock  func() time.Time
}

// New creates a Manager backed by the given JSON file path.
// Pass an empty path for an in-memory only manager.
func New(path string, clock func() time.Time) (*Manager, error) {
	if clock == nil {
		clock = time.Now
	}
	m := &Manager{path: path, clock: clock, leases: make(map[string]Lease)}
	if path == "" {
		return m, nil
	}
	if err := m.load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return m, nil
}

// Acquire grants a lease for key to owner lasting duration d.
// Returns an error if the key is already leased and not expired.
func (m *Manager) Acquire(key, owner string, d time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if existing, ok := m.leases[key]; ok && !existing.IsExpired(m.clock()) {
		return errors.New("key already leased by " + existing.Owner)
	}
	m.leases[key] = Lease{Key: key, Owner: owner, ExpiresAt: m.clock().Add(d)}
	return m.save()
}

// Release removes the lease for key if owned by owner.
func (m *Manager) Release(key, owner string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	l, ok := m.leases[key]
	if !ok {
		return nil
	}
	if l.Owner != owner {
		return errors.New("lease owned by " + l.Owner)
	}
	delete(m.leases, key)
	return m.save()
}

// Get returns the lease for key, if present.
func (m *Manager) Get(key string) (Lease, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	l, ok := m.leases[key]
	return l, ok
}

// PurgeExpired removes all expired leases and persists the result.
func (m *Manager) PurgeExpired() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := m.clock()
	for k, l := range m.leases {
		if l.IsExpired(now) {
			delete(m.leases, k)
		}
	}
	return m.save()
}

func (m *Manager) load() error {
	data, err := os.ReadFile(m.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &m.leases)
}

func (m *Manager) save() error {
	if m.path == "" {
		return nil
	}
	data, err := json.MarshalIndent(m.leases, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.path, data, 0o600)
}
