package main

import (
	"log"
	"time"

	"github.com/ucaptcha/backend-go/challenge"
	"github.com/ucaptcha/backend-go/config"
	"github.com/ucaptcha/backend-go/keys"
	"github.com/ucaptcha/backend-go/server"
)

func main() {
	if err := config.LoadConfig("config.example.yaml"); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	initialKey, err := keys.GenerateRSAKey(config.GlobalConfig.KeyLength)
	if err != nil {
		log.Fatalf("Failed to generate initial key: %v", err)
	}
	keys.AddKey(initialKey)

	challenge.InitializeStorage()

	go func() {
		interval := config.GlobalConfig.KeyRotationInterval
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			log.Println("Rotating RSA keys...")
			newKey, err := keys.GenerateRSAKey(config.GlobalConfig.KeyLength)
			if err != nil {
				log.Printf("Failed to generate new key: %v", err)
				continue
			}
			keys.AddKey(newKey)
			log.Println("RSA keys rotated.")
			// In a real system, you might need to handle old keys for ongoing challenges.
		}
	}()

	router := server.SetupRouter()
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
