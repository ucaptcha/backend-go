package types

import (
	"math/big"
	"time"
)

// Challenge represents the data associated with a cryptographic challenge.
type Challenge struct {
	ID        string
	G         *big.Int
	N         *big.Int // N is still needed for generating g, but P/Q are not
	T         int64    // Difficulty (number of iterations)
	CreatedAt time.Time
	KeyID     string // Reference to the key used for this challenge
}
