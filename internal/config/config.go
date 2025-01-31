package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI            string
	MongoDBName         string
	MongoCollectionName string
	Port                string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	config := &Config{
		MongoURI:            os.Getenv("MONGODB_URI"),
		MongoDBName:         os.Getenv("MONGODB_DB_NAME"),
		MongoCollectionName: os.Getenv("MONGODB_COLLECTION_NAME"),
		Port:                os.Getenv("PORT"),
	}

	// Validate required environment variables
	if config.MongoURI == "" {
		return nil, fmt.Errorf("MONGODB_URI is required")
	}
	if config.MongoDBName == "" {
		return nil, fmt.Errorf("MONGODB_DB_NAME is required")
	}
	if config.MongoCollectionName == "" {
		return nil, fmt.Errorf("MONGODB_COLLECTION_NAME is required")
	}
	if config.Port == "" {
		config.Port = "8080" // Default port if not specified
	}

	return config, nil
}
