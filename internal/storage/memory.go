package storage

import (
	"sync"

	"github.com/perceptumx/percepta/internal/core"
)

// MemoryStorage is an in-memory storage stub for MVP
// Observations stored in slice, not persisted
type MemoryStorage struct {
	mu           sync.RWMutex
	observations []core.Observation
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		observations: make([]core.Observation, 0),
	}
}

func (m *MemoryStorage) Save(obs core.Observation) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.observations = append(m.observations, obs)
	return nil
}

func (m *MemoryStorage) Query(deviceID string, limit int) ([]core.Observation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var filtered []core.Observation
	for _, obs := range m.observations {
		if deviceID == "" || obs.DeviceID == deviceID {
			filtered = append(filtered, obs)
		}
	}

	// Apply limit
	if limit > 0 && len(filtered) > limit {
		filtered = filtered[len(filtered)-limit:]
	}

	return filtered, nil
}

func (m *MemoryStorage) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.observations)
}
