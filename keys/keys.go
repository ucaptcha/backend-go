package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"log"
	"math/big"
	mrand "math/rand"
	"sync"
	"time"

	"github.com/ucaptcha/backend-go/config"
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
	if len(activeKeys) == 0 {
		key, err := GenerateRSAKey(config.GlobalConfig.KeyLength)
		if err != nil {
			return nil
		}
		AddKey(key)
		log.Println("No key found. Generated a new key.")
		return key
	}
	randomIndex := mrand.Intn(len(activeKeys))
	return activeKeys[randomIndex]
}

func RemoveOldKey() {
	keyMutex.Lock()
	defer keyMutex.Unlock()
	if len(activeKeys) > 0 {
		activeKeys = activeKeys[1:]
	}
}
