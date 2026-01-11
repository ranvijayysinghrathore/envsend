package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yourusername/envsend/backend/internal/api"
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

	log.Printf("Starting EnvSend API server in %s mode...", cfg.Server.Env)

	// Initialize storage layers
	postgres, err := storage.NewPostgresRepository(cfg.Database.URL)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer postgres.Close()
	log.Println("✓ Connected to PostgreSQL")

	redis, err := storage.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redis.Close()
	log.Println("✓ Connected to Redis")

	s3, err := storage.NewS3Storage(cfg.S3)
	if err != nil {
		log.Fatalf("Failed to connect to MinIO/S3: %v", err)
	}
	log.Println("✓ Connected to MinIO/S3")

	// Initialize services
	secretService := services.NewSecretService(postgres, s3, redis)
	log.Println("✓ Initialized services")

	// Create router
	router := api.NewRouter(cfg, secretService, redis)
	log.Println("✓ Configured router")

	// Create HTTP server
	server := &http.Server{
		Addr:         cfg.Server.Address(),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("🚀 Server listening on %s", cfg.Server.Address())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
