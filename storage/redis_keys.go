package storage

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/ucaptcha/backend-go/config"
)

// RedisKeyStorage is a Redis implementation of the KeyStorage interface.
type RedisKeyStorage struct {
	client *redis.Client
	prefix string // Prefix for Redis keys to avoid collisions
}

// NewRedisKeyStorage creates a new RedisKeyStorage instance.
func NewRedisKeyStorage(cfg config.RedisConfig) KeyStorage {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	return &RedisKeyStorage{
		client: client,
		prefix: "ucaptcha:key:", // Define a prefix for key storage
	}
}

func (s *RedisKeyStorage) GetKeyCount() (int, error) {
	ctx := s.client.Context()
	count, err := s.client.Keys(ctx, s.prefix+"*").Result()
	if err != nil {
		return 0, err
	}
	return len(count), nil
}

// SaveKey stores a key pair in Redis.
func (s *RedisKeyStorage) SaveKey(key *KeyPair) error {
	ctx := s.client.Context()
	redisKey := s.prefix + key.ID
	if key.ID == "" {
		return fmt.Errorf("key must have an ID")
	}

	// Serialize the KeyPair struct to JSON
	jsonData, err := json.Marshal(key)
	if err != nil {
		return fmt.Errorf("failed to marshal key pair: %v", err)
	}

	// Store the JSON string in Redis
	err = s.client.Set(ctx, redisKey, jsonData, 0).Err() // 0 means no expiration
	if err != nil {
		return fmt.Errorf("failed to save key to Redis: %v", err)
	}
	return nil
}

// GetRandomKey retrieves a random key pair from Redis.
func (s *RedisKeyStorage) GetRandomKey() (*KeyPair, error) {
	ctx := s.client.Context()

	// Get all keys with our prefix
	keys, err := s.client.Keys(ctx, s.prefix+"*").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get keys: %v", err)
	}
	if len(keys) == 0 {
		return nil, nil // No keys exist
	}

	// Select first key (simple approach - could be made truly random if needed)
	keyID := keys[0]

	// Get the key data
	jsonData, err := s.client.Get(ctx, keyID).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get key data: %v", err)
	}

	var key KeyPair
	err = json.Unmarshal([]byte(jsonData), &key)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal key pair: %v", err)
	}

	return &key, nil
}

// GetKey retrieves a key pair from Redis by its ID.
func (s *RedisKeyStorage) GetKey(id string) (*KeyPair, error) {
	ctx := s.client.Context()
	redisKey := s.prefix + id

	jsonData, err := s.client.Get(ctx, redisKey).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("key not found: %s", id)
	} else if err != nil {
		return nil, fmt.Errorf("failed to get key from Redis: %v", err)
	}

	var key KeyPair
	err = json.Unmarshal([]byte(jsonData), &key)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal key pair: %v", err)
	}

	return &key, nil
}

// DeleteKey removes a key pair from Redis by its ID.
func (s *RedisKeyStorage) DeleteKey(id string) error {
	ctx := s.client.Context()
	redisKey := s.prefix + id
	return s.client.Del(ctx, redisKey).Err()
}

// GetAllKeys retrieves all key pairs currently stored in Redis.
// Note: This can be inefficient in Redis with many keys. Consider alternatives if performance is critical.
func (s *RedisKeyStorage) GetAllKeys() ([]*KeyPair, error) {
	ctx := s.client.Context()
	var keyList []*KeyPair

	iter := s.client.Scan(ctx, 0, s.prefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		redisKey := iter.Val()
		jsonData, err := s.client.Get(ctx, redisKey).Result()
		if err != nil {
			// Log error or handle differently? For now, continue to next key.
			fmt.Printf("Error getting key %s: %v\n", redisKey, err)
			continue
		}

		var key KeyPair
		err = json.Unmarshal([]byte(jsonData), &key)
		if err != nil {
			fmt.Printf("Error unmarshalling key %s: %v\n", redisKey, err)
			continue
		}
		keyList = append(keyList, &key)
	}
	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("error iterating keys in Redis: %v", err)
	}

	return keyList, nil
}
