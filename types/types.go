package types

import (
	"math/big"
	"time"
)

// Challenge represents the data associated with a cryptographic challenge.
type Challenge struct {
	ID        string
	G         *big.Int
	N         *big.Int
	T         int64 // Difficulty (number of iterations)
	CreatedAt time.Time
	P         *big.Int // Store p and q for verification
	Q         *big.Int
}
