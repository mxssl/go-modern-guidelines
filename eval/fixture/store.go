package legacy

import (
	"context"
	"errors"
	"sync"
)

// Store is an in-memory key/value store that must be flushed on shutdown.
type Store struct {
	mu      sync.Mutex
	data    map[string]string
	flushed bool
}

// NewStore returns an empty store.
func NewStore() *Store {
	return &Store{data: make(map[string]string)}
}

// Put stores a value.
func (s *Store) Put(ctx context.Context, k, v string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[k] = v
	return nil
}

// Get reads a value.
func (s *Store) Get(ctx context.Context, k string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.data[k]
	if !ok {
		return "", errors.New("not found")
	}
	return v, nil
}

// Flush persists pending writes. It fails if the context is already done.
func (s *Store) Flush(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.flushed = true
	return nil
}

// Flushed reports whether Flush completed.
func (s *Store) Flushed() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.flushed
}
