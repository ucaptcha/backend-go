package challenge

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ucaptcha/backend-go/config"
	"github.com/ucaptcha/backend-go/keys"
	"github.com/ucaptcha/backend-go/storage"
	"github.com/ucaptcha/backend-go/types"
)

var globalStorage storage.Storage

// InitializeStorage initializes the global storage based on the configuration.
func InitializeStorage() {
	if config.GlobalConfig.Mode == "redis" {
		// Assuming NewRedisStorage is now in the storage package
		globalStorage = storage.NewRedisStorage(config.GlobalConfig.Redis)
	} else {
		// Assuming NewMemoryStorage is now in the storage package
		globalStorage = storage.NewMemoryStorage()
	}
}

// NewChallenge creates and stores a new challenge.
func NewChallenge() (*types.Challenge, error) {
	keyPair := keys.GetActiveKey()
	if keyPair == nil {
		return nil, fmt.Errorf("no active keys available")
	}

	challengeID := GenerateRandomID()
	g := GenerateValidG(keyPair.Components.N)
	difficulty := int64(1000000) // Example difficulty

	challenge := &types.Challenge{
		ID:        challengeID,
		G:         g,
		N:         keyPair.Components.N,
		T:         difficulty,
		CreatedAt: time.Now(),
		P:         keyPair.Components.P,
		Q:         keyPair.Components.Q,
	}

	if err := globalStorage.Save(challenge); err != nil {
		return nil, err
	}

	return challenge, nil
}

// GetChallenge retrieves a challenge by its ID.
func GetChallenge(id string) (*types.Challenge, error) {
	return globalStorage.Get(id)
}

// VerifyChallenge verifies the provided solution against the stored challenge.
func VerifyChallenge(id string, yStr string) int8 {
	challenge, err := globalStorage.Get(id)
	if err != nil {
		return 2 // Challenge not found
	}
	defer globalStorage.Delete(id) // Delete the challenge after verification

	y := new(big.Int)
	y, ok := y.SetString(yStr, 10)
	if !ok {
		return 3 // Invalid y format
	}

	pPrime := new(big.Int).Div(new(big.Int).Sub(challenge.P, big.NewInt(1)), big.NewInt(2))
	qPrime := new(big.Int).Div(new(big.Int).Sub(challenge.Q, big.NewInt(1)), big.NewInt(2))

	// Parallel computation of exponents modulo p and modulo q
	eP := new(big.Int).Exp(big.NewInt(2), big.NewInt(challenge.T), pPrime)
	eQ := new(big.Int).Exp(big.NewInt(2), big.NewInt(challenge.T), qPrime)

	// Decompose base and result to modulo p / modulo q
	gP := new(big.Int).Mod(challenge.G, challenge.P)
	gQ := new(big.Int).Mod(challenge.G, challenge.Q)
	yP := new(big.Int).Exp(gP, eP, challenge.P)
	yQ := new(big.Int).Exp(gQ, eQ, challenge.Q)

	// Directly verify the modulus decomposition result
	yp := new(big.Int).Mod(y, challenge.P)
	yq := new(big.Int).Mod(y, challenge.Q)

	if yp.Cmp(yP) == 0 && yq.Cmp(yQ) == 0 {
		return 1 // Verification successful
	}
	return 0 // Verification failed
}
