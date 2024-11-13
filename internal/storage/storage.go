package storage

import (
	"errors"
	"sync"
	"time"
)

// Storage struct represents an in-memory key-value store
type Storage struct {
	data       map[string]string
	hset       map[string]map[string]string
	expiration map[string]time.Time
	mu         sync.RWMutex
}

// NewStorage initializes a new Storage instance
func NewStorage() *Storage {
	store := &Storage{
		data:       make(map[string]string),
		hset:       make(map[string]map[string]string),
		expiration: make(map[string]time.Time),
	}

	// Start a goroutine to periodically remove expired keys
	go store.startExpirationCheck()

	return store
}

// Set stores a key-value pair with optional TTL
func (s *Storage) Set(key, value string, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value
	if ttl > 0 {
		s.expiration[key] = time.Now().Add(ttl)
	} else {
		delete(s.expiration, key) // remove any existing expiration
	}
}

// HSet stores a key-value pair for a given hash in hset
func (s *Storage) HSet(hash, key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.hset[hash]; !ok {
		s.hset[hash] = map[string]string{}
	}

	s.hset[hash][key] = value
}

// Get retrieves the value for a given key, considering TTL
func (s *Storage) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Check if the key has expired
	if expTime, exists := s.expiration[key]; exists && time.Now().After(expTime) {
		return "", false
	}

	value, found := s.data[key]
	return value, found
}

// HGet retrieves the value for a given hash and key in hset
func (s *Storage) HGet(hash, key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, found := s.hset[hash][key]
	return value, found
}

// HGetAll retrieves all values for a given hash in hset
func (s *Storage) HGetAll(hash string) (map[string]string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, found := s.hset[hash]
	return value, found
}

// Delete a key from storage
func (s *Storage) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, found := s.data[key]
	if !found {
		return errors.New("key not found")
	}

	delete(s.data, key)
	return nil
}

// HDelete deletes a key for a given hash in hset
func (s *Storage) HDelete(hash, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, found := s.hset[hash][key]
	if !found {
		return errors.New("key not found")
	}

	delete(s.hset[hash], key)
	return nil
}

// HDeleteAll deletes all values for a given hash in hset
func (s *Storage) HDeleteAll(hash string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, found := s.hset[hash]
	if !found {
		return errors.New("key not found")
	}

	delete(s.hset, hash)
	return nil
}

// startExpirationCheck periodically removes expired keys
func (s *Storage) startExpirationCheck() {
	for {
		time.Sleep(time.Second)

		s.mu.Lock()
		now := time.Now()

		for key, expTime := range s.expiration {
			if now.After(expTime) {
				delete(s.data, key)
				delete(s.expiration, key)
			}
		}

		s.mu.Unlock()
	}
}
