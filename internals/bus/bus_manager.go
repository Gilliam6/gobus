package bus

import (
	"fmt"
	"sync"
)

type Manager[T any] struct {
	Buses map[string]*Bus[T]
	mu    sync.RWMutex
}

func NewManager[T any]() *Manager[T] {
	return &Manager[T]{
		Buses: make(map[string]*Bus[T]),
	}
}

func (m *Manager[T]) GetBusSet() map[string]bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	set := make(map[string]bool)
	for id, _ := range m.Buses {
		set[id] = true
	}
	return set
}

func (m *Manager[T]) Subscribe(id string) (chan T, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if bus, ok := m.Buses[id]; !ok {
		return nil, fmt.Errorf("bus %s not found", id)
	} else {
		return bus.Subscribe(), nil
	}
}

func (m *Manager[T]) Unsubscribe(id string, ch chan T) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if bus, ok := m.Buses[id]; !ok {
		return fmt.Errorf("bus %s not found", id)
	} else {
		return bus.Unsubscribe(ch)
	}
}

func (m *Manager[T]) AddBus(id string) (*Bus[T], error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.Buses[id]; ok {
		return nil, fmt.Errorf("bus %s already exists", id)
	}
	bus := NewBus[T](1000)
	m.Buses[id] = bus
	return bus, nil
}

func (m *Manager[T]) DeleteBus(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.Buses[id]; !ok {
		return fmt.Errorf("bus %s not found", id)
	}
	m.Buses[id].Close()
	delete(m.Buses, id)
	return nil
}
