package storage

import (
	"fmt"
	"sync"
)

// MemoryKeyStorage is an in-memory implementation of the KeyStorage interface.
type MemoryKeyStorage struct {
	keys map[string]*KeyPair
	mu   sync.RWMutex
}

// NewMemoryKeyStorage creates a new MemoryKeyStorage instance.
func NewMemoryKeyStorage() KeyStorage {
	return &MemoryKeyStorage{
		keys: make(map[string]*KeyPair),
	}
}

// SaveKey stores a key pair in memory.
func (s *MemoryKeyStorage) SaveKey(key *KeyPair) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if key.ID == "" {
		return fmt.Errorf("key must have an ID")
	}
	s.keys[key.ID] = key
	return nil
}

// GetKey retrieves a key pair from memory by its ID.
func (s *MemoryKeyStorage) GetKey(id string) (*KeyPair, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	key, ok := s.keys[id]
	if !ok {
		return nil, fmt.Errorf("key not found: %s", id)
	}
	return key, nil
}

// DeleteKey removes a key pair from memory by its ID.
func (s *MemoryKeyStorage) DeleteKey(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.keys, id)
	return nil
}

// GetAllKeys retrieves all key pairs currently stored in memory.
func (s *MemoryKeyStorage) GetAllKeys() ([]*KeyPair, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	keyList := make([]*KeyPair, 0, len(s.keys))
	for _, key := range s.keys {
		keyList = append(keyList, key)
	}
	return keyList, nil
}

func (s *MemoryKeyStorage) GetKeyCount() (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.keys), nil
}
