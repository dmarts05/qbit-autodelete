package main

import (
	"log"

	"github.com/dmarts05/qbit-autodelete/internal/config"
	"github.com/dmarts05/qbit-autodelete/internal/qbitmanager"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	manager, err := qbitmanager.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create qBittorrent manager: %v", err)
	}
	manager.Run()
}
