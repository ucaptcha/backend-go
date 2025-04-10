package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"math/big"
	"sync"
	"time"
)

type RSAComponents struct {
	P *big.Int
	Q *big.Int
	N *big.Int
}

type KeyPair struct {
	Components  RSAComponents
	GeneratedAt time.Time
}

var activeKeys []*KeyPair
var keyMutex sync.RWMutex

func GenerateRSAKey(keyLength int) (*KeyPair, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, keyLength)
	if err != nil {
		return nil, err
	}

	return &KeyPair{
		Components: RSAComponents{
			P: privateKey.Primes[0],
			Q: privateKey.Primes[1],
			N: privateKey.N,
		},
		GeneratedAt: time.Now(),
	}, nil
}

func AddKey(key *KeyPair) {
	keyMutex.Lock()
	defer keyMutex.Unlock()
	activeKeys = append(activeKeys, key)
}

func GetActiveKey() *KeyPair {
	keyMutex.RLock()
	defer keyMutex.RUnlock()
	if len(activeKeys) > 0 {
		return activeKeys[0] // For now, just return the first active key
	}
	return nil
}

// In a real scenario, you'd have a background process to rotate keys
// based on the configuration. This is a simplified example.
