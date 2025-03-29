package bus

import (
	"fmt"
	"sync"
)

type Manager struct {
	Buses map[string]*Bus[[]byte]
	mu    sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		Buses: make(map[string]*Bus[[]byte]),
	}
}

func (m *Manager) GetBusSet() map[string]bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	set := make(map[string]bool)
	for id, _ := range m.Buses {
		set[id] = true
	}
	return set
}

func (m *Manager) Subscribe(id string) (chan []byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if bus, ok := m.Buses[id]; !ok {
		return nil, fmt.Errorf("bus %s not found", id)
	} else {
		return bus.Subscribe(), nil
	}
}

func (m *Manager) Unsubscribe(id string, ch chan []byte) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if bus, ok := m.Buses[id]; !ok {
		return fmt.Errorf("bus %s not found", id)
	} else {
		return bus.Unsubscribe(ch)
	}
}

func (m *Manager) AddBus(id string) (*Bus[[]byte], error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.Buses[id]; ok {
		return nil, fmt.Errorf("bus %s already exists", id)
	}
	bus := NewBus[[]byte](1000)
	m.Buses[id] = bus
	return bus, nil
}

func (m *Manager) DeleteBus(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.Buses[id]; !ok {
		return fmt.Errorf("bus %s not found", id)
	}
	m.Buses[id].Close()
	delete(m.Buses, id)
	return nil
}
