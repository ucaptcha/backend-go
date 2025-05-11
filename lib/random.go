package lib

import (
	"crypto/rand"
	"math/big"
	"strings"
)

// GenerateRandomBigInt generates a random big.Int within the specified range [min, max].
func GenerateRandomBigInt(min, max *big.Int) *big.Int {
	rangeVal := new(big.Int).Sub(max, min)
	bitLength := rangeVal.BitLen()
	byteLength := (bitLength + 7) / 8
	mask := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), uint(bitLength)), big.NewInt(1))
	result := new(big.Int)

	for {
		randomBytes := make([]byte, byteLength)
		_, err := rand.Read(randomBytes)
		if err != nil {
			panic(err) // Handle error appropriately
		}
		result.SetBytes(randomBytes)
		result.And(result, mask) // Ensure it does not exceed bitLength bits
		if result.Cmp(rangeVal) <= 0 {
			break
		}
	}
	return new(big.Int).Add(min, result)
}

// GenerateValidG generates a valid 'g' value based on the provided 'N'.
func GenerateValidG(N *big.Int) *big.Int {
	if N.Cmp(big.NewInt(4)) <= 0 {
		panic("N must be > 4")
	}
	one := big.NewInt(1)
	nMinusOne := new(big.Int).Sub(N, one)
	zero := big.NewInt(0)

	for {
		r := GenerateRandomBigInt(big.NewInt(2), nMinusOne)
		g := new(big.Int).Exp(r, big.NewInt(2), N)
		if g.Cmp(one) != 0 && g.Cmp(zero) != 0 && g.Cmp(nMinusOne) != 0 {
			return g
		}
	}
}

// GenerateRandomID generates a unique random ID.
func GenerateRandomID() string {
	const length = 10
	const allowedChars = "abcdefghijkmnpqrstuvwxyz23456789ABCDEFGHJKLMNPQRSTUVWXYZ"
	const allowedCharsLen = len(allowedChars)
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	var sb strings.Builder
	for _, byteVal := range b {
		sb.WriteByte(allowedChars[int(byteVal)%allowedCharsLen])
	}
	return sb.String()
}
