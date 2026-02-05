package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort      string
	DBUrl           string
	RustFSEndpoint  string
	RustFSAccessKey string
	RustFSSecretKey string
	RustFSBucket    string
	RustFSRegion    string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	return &Config{
		ServerPort:      getEnv("SERVER_PORT", ":8080"),
		DBUrl:           getEnv("DB_URL", ""),
		RustFSEndpoint:  getEnv("RUSTFS_ENDPOINT", "http://localhost:9000"),
		RustFSAccessKey: getEnv("RUSTFS_ACCESS_KEY", ""),
		RustFSSecretKey: getEnv("RUSTFS_SECRET_KEY", ""),
		RustFSBucket:    getEnv("RUSTFS_BUCKET", "documents"),
		RustFSRegion:    getEnv("RUSTFS_REGION", "us-east-1"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
