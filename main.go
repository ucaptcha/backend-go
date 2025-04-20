package main

import (
	"log"
	"strconv"
	"time"

	"github.com/ucaptcha/backend-go/challenge"
	"github.com/ucaptcha/backend-go/config"
	"github.com/ucaptcha/backend-go/keys"
	"github.com/ucaptcha/backend-go/server"
	"github.com/ucaptcha/backend-go/storage"
)

func main() {
	if err := config.LoadConfig("config.yaml"); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	keyPoolSize := config.GlobalConfig.KeyPoolSize

	// Initialize storage based on config
	var keyStorage storage.KeyStorage
	var challengeStorage storage.ChallengeStorage

	if config.GlobalConfig.KeysStorage == "redis" {
		keyStorage = storage.NewRedisKeyStorage(config.GlobalConfig.Redis)
	} else {
		keyStorage = storage.NewMemoryKeyStorage()
	}
	if config.GlobalConfig.ChallengeStorage == "redis" {
		challengeStorage = storage.NewRedisChallengeStorage(config.GlobalConfig.Redis)
	} else {
		challengeStorage = storage.NewMemoryChallengeStorage()
	}

	keyManager := keys.NewKeyManager(keyStorage, config.GlobalConfig.KeyLength)

	// Initialize challenge package
	challenge.InitializeStorage(challengeStorage, keyManager)

	currentKeyCount, err := keyStorage.GetKeyCount()
	if err != nil {
		log.Fatalf("Failed to get key count: %v", err)
	}

	// Generate initial keys if needed
	if currentKeyCount < keyPoolSize {
		for range keyPoolSize - currentKeyCount {
			_, err := keyManager.AddKey()
			if err != nil {
				log.Fatalf("Failed to generate initial key: %v", err)
			}
		}
		log.Printf("Generated %d initial keys", keyPoolSize-currentKeyCount)
	}

	currentKeyCount, err = keyStorage.GetKeyCount()
	if err != nil {
		log.Fatalf("Failed to get key count: %v", err)
	}

	log.Printf("Current key pool size: %d", currentKeyCount)

	go func() {
		interval := config.GlobalConfig.KeyRotationInterval
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			log.Println("Generating a new RSA key...")
			allKeys, err := keyStorage.GetAllKeys()
			if err != nil {
				log.Printf("Failed to get all keys: %v", err)
				continue
			}

			// Add new key
			_, err = keyManager.AddKey()
			if err != nil {
				log.Printf("Failed to generate new key: %v", err)
				continue
			}
			log.Println("RSA key generated.")

			// Remove oldest key if we have more than one key
			if len(allKeys) > 0 {
				oldestKey := allKeys[0]
				for _, key := range allKeys[1:] {
					if key.GeneratedAt.Before(oldestKey.GeneratedAt) {
						oldestKey = key
					}
				}
				if err := keyManager.RemoveKey(oldestKey.ID); err != nil {
					log.Printf("Failed to remove old key: %v", err)
				} else {
					log.Println("Removed old key.")
				}
			}
		}
	}()

	router := server.SetupRouter()
	if err := router.Run(config.GlobalConfig.Host + ":" + strconv.Itoa(config.GlobalConfig.Port)); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
