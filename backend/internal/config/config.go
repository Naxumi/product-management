package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	App      AppConfig
	Storage  StorageConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

// AppConfig holds application configuration
type AppConfig struct {
	Port     int
	Env      string
	LogLevel string
}

type StorageConfig struct {
	Type     string // "local", "minio", "s3"
	BasePath string // "./storage"
	BaseURL  string // "http://localhost:8080/uploads"

	// MinIO/S3 config (for future)
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := &Config{}

	// Database configuration
	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	config.Database = DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     dbPort,
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", ""),
		Name:     getEnv("DB_NAME", "product_management"),
		SSLMode:  getEnv("DB_SSL_MODE", "disable"),
	}

	// Application configuration
	appPort, err := strconv.Atoi(getEnv("APP_PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid APP_PORT: %w", err)
	}

	config.App = AppConfig{
		Port:     appPort,
		Env:      getEnv("APP_ENV", "development"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}
	// Storage Configuration
	storageType := getEnv("STORAGE_TYPE", "")
	basePath := getEnv("BASE_PATH", "")
	baseURL := getEnv("BASE_URL", "")
	config.Storage = StorageConfig{
		Type:     storageType,
		BasePath: basePath,
		BaseURL:  baseURL,
	}

	// Validate required fields
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Database.Password == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}
	if c.Storage.Type == "" {
		return fmt.Errorf("STORAGE_TYPE is required")
	}
	if c.Storage.BasePath == "" {
		return fmt.Errorf("BASE_PATH is required")
	}
	if c.Storage.BaseURL == "" {
		return fmt.Errorf("BASE_URL is required")
	}
	return nil
}

// DatabaseURL returns the PostgreSQL connection string
func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
