package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	S3       S3Config
	Security SecurityConfig
	Worker   WorkerConfig
	CORS     CORSConfig
	Logging  LoggingConfig
}

// ServerConfig holds HTTP server configuration.
type ServerConfig struct {
	Host string
	Port int
	Env  string
}

// DatabaseConfig holds PostgreSQL configuration.
type DatabaseConfig struct {
	URL            string
	MaxConnections int
	MaxIdle        int
}

// RedisConfig holds Redis configuration.
type RedisConfig struct {
	URL      string
	Password string
}

// S3Config holds MinIO/S3 configuration.
type S3Config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
	Region    string
}

// SecurityConfig holds security-related configuration.
type SecurityConfig struct {
	RateLimitRequests int
	RateLimitWindow   time.Duration
	MaxSecretSize     int64
	DefaultExpiry     time.Duration
	MaxExpiry         time.Duration
}

// WorkerConfig holds background worker configuration.
type WorkerConfig struct {
	CleanupInterval time.Duration
	BatchSize       int
}

// CORSConfig holds CORS configuration.
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

// LoggingConfig holds logging configuration.
type LoggingConfig struct {
	Level  string
	Format string
}

// Load loads configuration from environment variables.
func Load() (*Config, error) {
	// Load .env file if it exists (ignore errors in production)
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnvAsInt("SERVER_PORT", 8080),
			Env:  getEnv("SERVER_ENV", "development"),
		},
		Database: DatabaseConfig{
			URL:            getEnv("DATABASE_URL", "postgres://envsend:envsend@localhost:5432/envsend?sslmode=disable"),
			MaxConnections: getEnvAsInt("DATABASE_MAX_CONNECTIONS", 25),
			MaxIdle:        getEnvAsInt("DATABASE_MAX_IDLE_CONNECTIONS", 5),
		},
		Redis: RedisConfig{
			URL:      getEnv("REDIS_URL", "redis://localhost:6379/0"),
			Password: getEnv("REDIS_PASSWORD", ""),
		},
		S3: S3Config{
			Endpoint:  getEnv("S3_ENDPOINT", "localhost:9000"),
			AccessKey: getEnv("S3_ACCESS_KEY", "minioadmin"),
			SecretKey: getEnv("S3_SECRET_KEY", "minioadmin"),
			Bucket:    getEnv("S3_BUCKET", "envsend-secrets"),
			UseSSL:    getEnvAsBool("S3_USE_SSL", false),
			Region:    getEnv("S3_REGION", "us-east-1"),
		},
		Security: SecurityConfig{
			RateLimitRequests: getEnvAsInt("RATE_LIMIT_REQUESTS", 100),
			RateLimitWindow:   getEnvAsDuration("RATE_LIMIT_WINDOW", 60*time.Second),
			MaxSecretSize:     getEnvAsInt64("MAX_SECRET_SIZE", 10*1024*1024), // 10 MB
			DefaultExpiry:     getEnvAsDuration("DEFAULT_SECRET_EXPIRY", 10*time.Minute),
			MaxExpiry:         getEnvAsDuration("MAX_SECRET_EXPIRY", 168*time.Hour), // 7 days
		},
		Worker: WorkerConfig{
			CleanupInterval: getEnvAsDuration("WORKER_CLEANUP_INTERVAL", 60*time.Second),
			BatchSize:       getEnvAsInt("WORKER_BATCH_SIZE", 100),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"*"}),
			AllowedMethods: getEnvAsSlice("CORS_ALLOWED_METHODS", []string{"GET", "POST", "DELETE", "OPTIONS"}),
			AllowedHeaders: getEnvAsSlice("CORS_ALLOWED_HEADERS", []string{"Content-Type", "Authorization"}),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}

	return cfg, nil
}

// Helper functions for environment variables

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	valueStr := os.Getenv(key)
	if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	// Simple comma-separated parsing
	var result []string
	for _, v := range splitAndTrim(valueStr, ",") {
		if v != "" {
			result = append(result, v)
		}
	}
	if len(result) == 0 {
		return defaultValue
	}
	return result
}

func splitAndTrim(s, sep string) []string {
	var result []string
	for _, part := range splitString(s, sep) {
		trimmed := trimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func splitString(s, sep string) []string {
	// Simple string split implementation
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	result = append(result, s[start:])
	return result
}

func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}

// Address returns the server address in host:port format.
func (c *ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// IsDevelopment returns true if running in development mode.
func (c *ServerConfig) IsDevelopment() bool {
	return c.Env == "development"
}

// IsProduction returns true if running in production mode.
func (c *ServerConfig) IsProduction() bool {
	return c.Env == "production"
}
