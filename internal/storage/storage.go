package storage

import (
	"errors"
	"sync"
)

// Storage struct represents an in-memory key-value store
type Storage struct {
	data map[string]string
	mu sync.RWMutex
}

// NewStorage initializes a new Storage instance
func NewStorage() *Storage {
	return &Storage{
		data: make(map[string]string),
	}
}

// Set stores a key-value pair
func (s *Storage) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

// Get retrieves the value for a given key
func (s *Storage) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, found := s.data[key]
	return value, found
}

// Delete a key from storage
func (s *Storage) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, found := s.data[key]
	if !found {
		return errors.New("Key not found")
	}
	delete(s.data, key)
	return nil
}