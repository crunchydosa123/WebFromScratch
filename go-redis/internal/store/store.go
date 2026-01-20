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

func (s *Store) Set(key, val string, ttlSecond *int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = val

	if ttlSecond != nil {
		s.expiry[key] = time.Now().Add(time.Duration(*ttlSecond) * time.Second)
	} else {
		delete(s.expiry, key)
	}
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	return val, ok
}
