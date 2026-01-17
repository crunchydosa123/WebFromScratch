package store

import (
	"sync"
	"time"
)

type Store struct {
	mu     sync.RWMutex
	data   map[string]string
	expiry map[string]time.Time
}

func New() *Store {
	return &Store{
		data:   make(map[string]string),
		expiry: make(map[string]time.Time),
	}
}

func (s *Store) Set(key, val string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = val
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	return val, ok
}
