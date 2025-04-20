package storage

import (
	"math/big"
	"time"

	"github.com/ucaptcha/backend-go/types"
)

// ChallengeStorage defines the interface for challenge storage operations.
type ChallengeStorage interface {
	Save(ch *types.Challenge) error
	Get(id string) (*types.Challenge, error)
	Delete(id string) error
}

// RSAComponents holds the components of an RSA key pair
type RSAComponents struct {
	P *big.Int `json:"p"`
	Q *big.Int `json:"q"`
	N *big.Int `json:"n"`
}

// KeyPair represents an RSA key pair with metadata
type KeyPair struct {
	ID          string        `json:"id"` // Unique identifier for the key pair
	Components  RSAComponents `json:"components"`
	GeneratedAt time.Time     `json:"generated_at"`
}

// KeyStorage defines the interface for key pair storage operations.
type KeyStorage interface {
	SaveKey(key *KeyPair) error
	GetKey(id string) (*KeyPair, error)
	DeleteKey(id string) error
	GetAllKeys() ([]*KeyPair, error)
	GetKeyCount() (int, error)
	GetRandomKey() (*KeyPair, error)
}
