package bus

import (
	"fmt"
	"sync"
)

type Bus[T any] struct {
	subscribers []chan T
	mutex       sync.RWMutex
	buffSize    int
}

func NewBus[T any](size int) *Bus[T] {
	return &Bus[T]{
		subscribers: make([]chan T, 0),
		buffSize:    size,
	}
}

func (b *Bus[T]) Subscribe() chan T {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	ch := make(chan T, b.buffSize)
	b.subscribers = append(b.subscribers, ch)
	return ch
}

func (b *Bus[T]) Unsubscribe(ch chan T) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	for i, c := range b.subscribers {
		if c == ch {
			b.subscribers = append(b.subscribers[:i], b.subscribers[i+1:]...)
			close(ch)
			return nil
		}
	}
	return fmt.Errorf("unsubscribing failed: channel does not exist")
}

func (b *Bus[T]) Publish(msg T) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	for _, c := range b.subscribers {
		select {
		case c <- msg:
		default:
		}
	}
}

func (b *Bus[T]) Close() {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	for _, ch := range b.subscribers {
		close(ch)
	}
	b.subscribers = nil
}
