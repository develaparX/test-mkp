package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Cache    CacheConfig    `json:"cache"`
	Logger   LoggerConfig   `json:"logger"`
	JWT      JWTConfig      `json:"jwt"`
}

type ServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
	SSLMode  string `json:"ssl_mode"`
}

type CacheConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type LoggerConfig struct {
	Level string `json:"level"`
	Type  string `json:"type"`
}

type JWTConfig struct {
	SecretKey string `json:"secret_key"`
	Issuer    string `json:"issuer"`
}

func LoadConfig(envPath string) (*Config, error) {

	if err := godotenv.Load(envPath); err != nil {

	}

	serverPort, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	if err != nil {
		serverPort = 8080
	}

	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5000"))
	if err != nil {
		dbPort = 5000
	}

	cachePort, err := strconv.Atoi(getEnv("CACHE_PORT", "6379"))
	if err != nil {
		cachePort = 6379
	}

	cacheDB, err := strconv.Atoi(getEnv("CACHE_DB", "0"))
	if err != nil {
		cacheDB = 0
	}

	config := &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "localhost"),
			Port: serverPort,
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     dbPort,
			Username: getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			Database: getEnv("DB_NAME", "belimang"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Cache: CacheConfig{
			Host:     getEnv("CACHE_HOST", "localhost"),
			Port:     cachePort,
			Password: getEnv("CACHE_PASSWORD", ""),
			DB:       cacheDB,
		},
		Logger: LoggerConfig{
			Level: getEnv("LOG_LEVEL", "info"),
			Type:  getEnv("LOG_TYPE", "simple"),
		},
		JWT: JWTConfig{
			SecretKey: getEnv("JWT_SECRET_KEY", "your-secret-key"),
			Issuer:    getEnv("JWT_ISSUER", "belimang-app"),
		},
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
