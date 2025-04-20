package storage

import (
	"fmt"
	"sync"

	"github.com/ucaptcha/backend-go/types"
)

// MemoryStorage is an in-memory implementation of the ChallengeStorage interface.
type MemoryStorage struct {
	challenges map[string]*types.Challenge
	mu         sync.Mutex
}

// NewMemoryChallengeStorage creates a new MemoryStorage instance.
func NewMemoryChallengeStorage() ChallengeStorage {
	return &MemoryStorage{
		challenges: make(map[string]*types.Challenge),
	}
}

// Save stores a challenge in memory.
func (s *MemoryStorage) Save(ch *types.Challenge) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.challenges[ch.ID] = ch
	return nil
}

// Get retrieves a challenge from memory by its ID.
func (s *MemoryStorage) Get(id string) (*types.Challenge, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	ch, ok := s.challenges[id]
	if !ok {
		return nil, fmt.Errorf("challenge not found: %s", id)
	}
	return ch, nil
}

// Delete removes a challenge from memory by its ID.
func (s *MemoryStorage) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.challenges, id)
	return nil
}
