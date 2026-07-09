package inmemorykvstore

import (
	"sync"

	"github.com/mirno/petshop/internal/usecases"
)

var _ usecases.Store[any] = (*InMemoryKVStore[any])(nil)

type InMemoryKVStore[T any] struct {
	data map[string]T
	mu   sync.RWMutex
}

func NewInMemoryKVStore[T any]() *InMemoryKVStore[T] {
	return &InMemoryKVStore[T]{
		data: make(map[string]T),
	}
}

func (store *InMemoryKVStore[T]) Get(key string) (T, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	value, exists := store.data[key]
	if !exists {
		return value, usecases.ErrKeyNotFoundError
	}

	return value, nil
}

func (store *InMemoryKVStore[T]) Keys() ([]string, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	var keys []string
	for key := range store.data {
		keys = append(keys, key)
	}

	return keys, nil
}

func (store *InMemoryKVStore[T]) Save(key string, value T) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	store.data[key] = value

	return nil
}

func (store *InMemoryKVStore[T]) Delete(key string) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	delete(store.data, key)

	return nil
}
