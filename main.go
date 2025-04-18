package main

import (
	"log"
	"strconv"
	"time"

	"github.com/ucaptcha/backend-go/challenge"
	"github.com/ucaptcha/backend-go/config"
	"github.com/ucaptcha/backend-go/keys"
	"github.com/ucaptcha/backend-go/server"
)

func main() {
	if err := config.LoadConfig("config.yaml"); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	keyPoolSize := config.GlobalConfig.KeyPoolSize

	// Generate initial keys
	for range keyPoolSize {
		initialKey, err := keys.GenerateRSAKey(config.GlobalConfig.KeyLength)
		if err != nil {
			log.Fatalf("Failed to generate initial key: %v", err)
		}
		keys.AddKey(initialKey)
	}

	log.Printf("Initialized a key pool with size %d", keyPoolSize)

	challenge.InitializeStorage()

	go func() {
		interval := config.GlobalConfig.KeyRotationInterval
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			log.Println("Generating a new RSA key...")
			newKey, err := keys.GenerateRSAKey(config.GlobalConfig.KeyLength)
			if err != nil {
				log.Printf("Failed to generate new key: %v", err)
				continue
			}
			keys.AddKey(newKey)
			log.Println("RSA keys generated.")
			keys.RemoveOldKey()
			log.Println("Removed an old key.")
		}
	}()

	router := server.SetupRouter()
	if err := router.Run(config.GlobalConfig.Host + ":" + strconv.Itoa(config.GlobalConfig.Port)); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
