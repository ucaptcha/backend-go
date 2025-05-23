package storage

import (
	"fmt"
	"math/big"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/ucaptcha/backend-go/config"
	"github.com/ucaptcha/backend-go/types"
)

// RedisStorage is a Redis implementation of the ChallengeStorage interface.
type RedisStorage struct {
	client *redis.Client
}

// NewRedisChallengeStorage creates a new RedisStorage instance.
func NewRedisChallengeStorage(cfg config.RedisConfig) ChallengeStorage {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	return &RedisStorage{client: client}
}

// Save stores a challenge in Redis.
func (s *RedisStorage) Save(ch *types.Challenge) error {
	ctx := s.client.Context()
	key := fmt.Sprintf("ucaptcha:challenge:%s", ch.ID)
	err := s.client.HSet(ctx, key,
		"id", ch.ID,
		"KeyID", ch.KeyID,
		"g", ch.G.String(),
		"n", ch.N.String(),
		"t", ch.T,
		"created_at", ch.CreatedAt.Format(time.RFC3339),
	).Err()
	if err != nil {
		return err
	}
	// Set an expiry for the challenge (e.g., 5 minutes)
	return s.client.Expire(ctx, key, 5*time.Minute).Err()
}

// Get retrieves a challenge from Redis by its ID.
func (s *RedisStorage) Get(id string) (*types.Challenge, error) {
	ctx := s.client.Context()
	key := fmt.Sprintf("ucaptcha:challenge:%s", id)
	result, err := s.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("challenge not found: %s", id)
	}

	g, _ := new(big.Int).SetString(result["g"], 10)
	n, _ := new(big.Int).SetString(result["n"], 10)
	t, _ := new(big.Int).SetString(result["t"], 10)
	createdAt, _ := time.Parse(time.RFC3339, result["created_at"])

	return &types.Challenge{
		ID:        id,
		G:         g,
		N:         n,
		T:         t.Int64(),
		CreatedAt: createdAt,
		KeyID:     result["KeyID"],
	}, nil
}

// Delete removes a challenge from Redis by its ID.
func (s *RedisStorage) Delete(id string) error {
	ctx := s.client.Context()
	key := fmt.Sprintf("ucaptcha:challenge:%s", id)
	return s.client.Del(ctx, key).Err()
}
