package storage

import (
	"github.com/ucaptcha/backend-go/types"
)

// Storage defines the interface for persisting and retrieving challenges.
type Storage interface {
	Save(challenge *types.Challenge) error
	Get(id string) (*types.Challenge, error)
	Delete(id string) error
}
