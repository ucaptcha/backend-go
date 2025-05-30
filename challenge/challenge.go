package challenge

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ucaptcha/backend-go/config"
	"github.com/ucaptcha/backend-go/keys"
	"github.com/ucaptcha/backend-go/lib"
	"github.com/ucaptcha/backend-go/storage"
	"github.com/ucaptcha/backend-go/types"
)

var (
	globalManager     *ChallengeManager
	globalManagerOnce sync.Once
)

// ChallengeManager handles the creation, retrieval, and verification of challenges.
type ChallengeManager struct {
	challengeStorage storage.ChallengeStorage
	keyManager       *keys.KeyManager
}

// NewChallengeManager creates a new ChallengeManager instance.
func NewChallengeManager(cs storage.ChallengeStorage, km *keys.KeyManager) *ChallengeManager {
	return &ChallengeManager{
		challengeStorage: cs,
		keyManager:       km,
	}
}

// InitializeStorage sets up the global ChallengeManager instance.
// Must be called before using NewChallenge() or VerifyChallenge().
func InitializeStorage(cs storage.ChallengeStorage, km *keys.KeyManager) {
	globalManagerOnce.Do(func() {
		globalManager = NewChallengeManager(cs, km)
	})
}

// NewChallenge creates a new challenge using the global manager.
func NewChallenge(difficulty ...int64) (*types.Challenge, error) {
	if globalManager == nil {
		return nil, fmt.Errorf("challenge storage not initialized")
	}
	return globalManager.NewChallenge(difficulty...)
}

// VerifyChallenge verifies a challenge using the global manager.
func VerifyChallenge(id string, yStr string) (int8, error) {
	if globalManager == nil {
		return 0, fmt.Errorf("challenge storage not initialized")
	}
	return globalManager.VerifyChallenge(id, yStr)
}

// NewChallenge creates and stores a new challenge.
func (cm *ChallengeManager) NewChallenge(difficulty ...int64) (*types.Challenge, error) {
	keyPair, err := cm.keyManager.GetRandomKey()

	if err != nil {
		return nil, fmt.Errorf("error getting random key: %v", err)
	}
	if keyPair == nil {
		// This case should ideally not happen if GetRandomKey guarantees a key
		return nil, fmt.Errorf("no active keys available, generation failed")
	}

	challengeID := lib.GenerateRandomID()
	// N is still needed for generating g, which is part of the public challenge
	g := lib.GenerateValidG(keyPair.Components.N)

	// Set difficulty to the provided value or use the default
	var diff int64 = config.GlobalConfig.Difficulty // Default difficulty
	if len(difficulty) > 0 {
		diff = difficulty[0]
	}

	challenge := &types.Challenge{
		ID:        challengeID,
		G:         g,
		N:         keyPair.Components.N, // N is public
		T:         diff,
		CreatedAt: time.Now(),
		KeyID:     keyPair.ID, // Store KeyID instead of P, Q
	}

	if err := cm.challengeStorage.Save(challenge); err != nil {
		return nil, fmt.Errorf("failed to save challenge: %v", err)
	}

	return challenge, nil
}

// GetChallenge retrieves a challenge by its ID.
func (cm *ChallengeManager) GetChallenge(id string) (*types.Challenge, error) {
	ch, err := cm.challengeStorage.Get(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get challenge %s: %v", id, err)
	}
	return ch, nil
}

// VerifyChallenge verifies the provided solution against the stored challenge.
func (cm *ChallengeManager) VerifyChallenge(id string, yStr string) (int8, error) {
	challenge, err := cm.challengeStorage.Get(id)
	if err != nil {
		return 2, fmt.Errorf("could not found challenge: %s", id) // Challenge not found
	}

	// Retrieve the key used for this challenge
	keyPair, err := cm.keyManager.GetKey(challenge.KeyID)
	if err != nil {
		return 4, fmt.Errorf("required key %s for challenge %s is missing, consider re-generating challenge", challenge.KeyID, id)
	}

	y := new(big.Int)
	y, ok := y.SetString(yStr, 10)
	if !ok {
		return 3, fmt.Errorf("invalid format for y: %s", yStr) // Invalid y format
	}

	// Perform verification using the retrieved key components
	pPrime := new(big.Int).Div(new(big.Int).Sub(keyPair.Components.P, big.NewInt(1)), big.NewInt(2))
	qPrime := new(big.Int).Div(new(big.Int).Sub(keyPair.Components.Q, big.NewInt(1)), big.NewInt(2))

	// Computation of exponents modulo p' and q'
	eP := new(big.Int).Exp(big.NewInt(2), big.NewInt(challenge.T), pPrime)
	eQ := new(big.Int).Exp(big.NewInt(2), big.NewInt(challenge.T), qPrime)

	// Decompose base and result to modulo p / modulo q
	gP := new(big.Int).Mod(challenge.G, keyPair.Components.P)
	gQ := new(big.Int).Mod(challenge.G, keyPair.Components.Q)
	yP := new(big.Int).Exp(gP, eP, keyPair.Components.P)
	yQ := new(big.Int).Exp(gQ, eQ, keyPair.Components.Q)

	// Directly verify the modulus decomposition result
	yp := new(big.Int).Mod(y, keyPair.Components.P)
	yq := new(big.Int).Mod(y, keyPair.Components.Q)

	if yp.Cmp(yP) == 0 && yq.Cmp(yQ) == 0 {
		// Verification successful, delete the challenge
		if delErr := cm.challengeStorage.Delete(id); delErr != nil {
			// Log the deletion error but still return success for verification
			fmt.Printf("Warning: Failed to delete challenge %s after successful verification: %v\n", id, delErr)
		}
		return 1, nil // Verification successful
	}

	// Verification failed, delete the challenge
	if delErr := cm.challengeStorage.Delete(id); delErr != nil {
		// Log the deletion error
		fmt.Printf("Warning: Failed to delete challenge %s after failed verification: %v\n", id, delErr)
	}
	return 0, nil // Verification failed
}
