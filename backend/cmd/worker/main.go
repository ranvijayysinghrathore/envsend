package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yourusername/envsend/backend/internal/config"
	"github.com/yourusername/envsend/backend/internal/services"
	"github.com/yourusername/envsend/backend/internal/storage"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Println("Starting EnvSend cleanup worker...")

	// Initialize storage layers
	postgres, err := storage.NewPostgresRepository(cfg.Database.URL)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer postgres.Close()
	log.Println("✓ Connected to PostgreSQL")

	s3, err := storage.NewS3Storage(cfg.S3)
	if err != nil {
		log.Fatalf("Failed to connect to MinIO/S3: %v", err)
	}
	log.Println("✓ Connected to MinIO/S3")

	redis, err := storage.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redis.Close()
	log.Println("✓ Connected to Redis")

	// Initialize services
	secretService := services.NewSecretService(postgres, s3, redis)
	log.Println("✓ Initialized services")

	// Create ticker for periodic cleanup
	ticker := time.NewTicker(cfg.Worker.CleanupInterval)
	defer ticker.Stop()

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("🧹 Worker started (cleanup interval: %s, batch size: %d)", cfg.Worker.CleanupInterval, cfg.Worker.BatchSize)

	// Run initial cleanup
	runCleanup(ctx, secretService, cfg.Worker.BatchSize)

	// Main worker loop
	for {
		select {
		case <-ticker.C:
			runCleanup(ctx, secretService, cfg.Worker.BatchSize)

		case <-quit:
			log.Println("Shutting down worker...")
			return
		}
	}
}

func runCleanup(ctx context.Context, service *services.SecretService, batchSize int) {
	log.Println("Running cleanup...")

	count, err := service.CleanupExpiredSecrets(ctx, batchSize)
	if err != nil {
		log.Printf("Cleanup error: %v", err)
		return
	}

	if count > 0 {
		log.Printf("✓ Cleaned up %d expired secrets", count)
	} else {
		log.Println("✓ No expired secrets to clean up")
	}
}
