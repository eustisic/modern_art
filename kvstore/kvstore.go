package kvstore

import "sync"

type StoreInterface interface {
	Insert(key, value string)
	Search(key string) (string, bool)
	Update(key, value string) bool
}

// Store represents a simple key-value store.
type Store struct {
	data map[string]string
	mu   sync.RWMutex // To ensure concurrent access
}

// NewStore initializes and returns a new Store.
func NewStore() *Store {
	return &Store{
		data: make(map[string]string),
	}
}

// Insert adds a new key-value pair to the store.
func (s *Store) Insert(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

// Search returns the value for a given key and a boolean indicating if the key was found.
func (s *Store) Search(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.data[key]
	return value, ok
}

// Update modifies the value of an existing key.
func (s *Store) Update(key, value string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[key]; ok {
		s.data[key] = value
		return true
	}
	return false
}
