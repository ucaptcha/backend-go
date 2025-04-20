package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
	mrand "math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/ucaptcha/backend-go/storage"
)

// KeyManager handles key generation, storage, and retrieval
type KeyManager struct {
	keyStorage storage.KeyStorage
	keyLength  int
	keyMutex   sync.RWMutex // Mutex for key generation/initialization logic
}

// NewKeyManager creates a new KeyManager instance.
func NewKeyManager(keyStorage storage.KeyStorage, keyLength int) *KeyManager {
	return &KeyManager{
		keyStorage: keyStorage,
		keyLength:  keyLength,
	}
}

// generateNewKey generates a new RSA key pair with a unique ID.
func generateNewKey(keyLength int) (*storage.KeyPair, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, keyLength)
	if err != nil {
		return nil, err
	}

	keyID := uuid.New().String()
	return &storage.KeyPair{
		ID: keyID,
		Components: storage.RSAComponents{
			P: privateKey.Primes[0],
			Q: privateKey.Primes[1],
			N: privateKey.N,
		},
		GeneratedAt: time.Now(),
	}, nil
}

// GetKey retrieves a key by its ID.
func (km *KeyManager) GetKey(id string) (*storage.KeyPair, error) {
	return km.keyStorage.GetKey(id)
}

// GetRandomKey retrieves a random key from storage.
// If no keys exist, it generates a new one, saves it, and returns it.
func (km *KeyManager) GetRandomKey() (*storage.KeyPair, error) {
	km.keyMutex.RLock() // Use RLock first for reading keys
	allKeys, err := km.keyStorage.GetAllKeys()
	km.keyMutex.RUnlock()

	if err != nil {
		return nil, fmt.Errorf("failed to get all keys: %v", err)
	}

	if len(allKeys) > 0 {
		mrand.Seed(time.Now().UnixNano()) // Ensure randomness
		randomIndex := mrand.Intn(len(allKeys))
		return allKeys[randomIndex], nil
	}

	// No keys found, need to generate one. Acquire write lock.
	km.keyMutex.Lock()
	defer km.keyMutex.Unlock()

	// Double-check if another goroutine generated a key while waiting for the lock
	allKeys, err = km.keyStorage.GetAllKeys()
	if err != nil {
		// Log this error, but proceed to generate a key as a fallback
		log.Printf("Warning: Failed to re-check keys after acquiring lock: %v", err)
	} else if len(allKeys) > 0 {
		// Another goroutine created a key, use that one
		mrand.Seed(time.Now().UnixNano())
		randomIndex := mrand.Intn(len(allKeys))
		return allKeys[randomIndex], nil
	}

	// Still no keys, generate, save, and return a new one
	log.Println("No keys found in storage. Generating a new key.")
	newKey, err := generateNewKey(km.keyLength)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new key: %v", err)
	}

	err = km.keyStorage.SaveKey(newKey)
	if err != nil {
		// Log the save error, but return the generated key anyway if needed elsewhere
		log.Printf("Warning: Failed to save newly generated key: %v", err)
		// Depending on requirements, you might want to return an error here instead
	}

	return newKey, nil
}

// AddKey generates a new key and saves it to storage.
func (km *KeyManager) AddKey() (*storage.KeyPair, error) {
	km.keyMutex.Lock() // Lock needed as it modifies storage
	defer km.keyMutex.Unlock()

	newKey, err := generateNewKey(km.keyLength)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new key: %v", err)
	}
	err = km.keyStorage.SaveKey(newKey)
	if err != nil {
		return nil, fmt.Errorf("failed to save new key: %v", err)
	}
	log.Printf("Added new key with ID: %s", newKey.ID)
	return newKey, nil
}

// RemoveKey removes a key by its ID.
func (km *KeyManager) RemoveKey(id string) error {
	km.keyMutex.Lock() // Lock needed as it modifies storage
	defer km.keyMutex.Unlock()

	err := km.keyStorage.DeleteKey(id)
	if err != nil {
		return fmt.Errorf("failed to delete key %s: %v", id, err)
	}
	log.Printf("Removed key with ID: %s", id)
	return nil
}

// Get key count
func (km *KeyManager) GetKeyCount() (int, error) {
	km.keyMutex.RLock() // Lock needed as it modifies storage
	defer km.keyMutex.RUnlock()
	keyCounts, err := km.keyStorage.GetKeyCount()
	if err != nil {
		return 0, err
	}
	return keyCounts, nil
}
