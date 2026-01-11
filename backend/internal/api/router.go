package api

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/yourusername/envsend/backend/internal/api/handlers"
	mw "github.com/yourusername/envsend/backend/internal/api/middleware"
	"github.com/yourusername/envsend/backend/internal/config"
	"github.com/yourusername/envsend/backend/internal/services"
	"github.com/yourusername/envsend/backend/internal/storage"
)

// NewRouter creates and configures the HTTP router.
func NewRouter(cfg *config.Config, secretService *services.SecretService, redis *storage.RedisClient) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(mw.SecurityHeaders)
	r.Use(mw.AuditLogger)

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CORS.AllowedOrigins,
		AllowedMethods:   cfg.CORS.AllowedMethods,
		AllowedHeaders:   cfg.CORS.AllowedHeaders,
		ExposedHeaders:   []string{"X-RateLimit-Limit", "X-RateLimit-Remaining"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Health check endpoints (no rate limiting)
	healthHandler := handlers.NewHealthHandler()
	r.Get("/health", healthHandler.Health)
	r.Get("/ready", healthHandler.Readiness)

	// API routes with rate limiting
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(mw.RateLimiter(redis, cfg.Security.RateLimitRequests, cfg.Security.RateLimitWindow))

		// Secrets endpoints
		secretsHandler := handlers.NewSecretsHandler(secretService)

		// POST /api/v1/secrets - Create secret
		r.With(mw.NoSecretLogging).Post("/secrets", secretsHandler.CreateSecret)

		// GET /api/v1/secrets/:id - Get secret
		r.With(mw.NoSecretLogging).Get("/secrets/{id}", secretsHandler.GetSecret)
	})

	// Short URL redirect (for user-friendly links)
	r.Get("/s/{id}", func(w http.ResponseWriter, r *http.Request) {
		secretID := chi.URLParam(r, "id")
		// TODO: Serve web UI for browser-based decryption
		// For now, return JSON
		w.Header().Set("Location", "/api/v1/secrets/"+secretID)
		w.WriteHeader(http.StatusTemporaryRedirect)
	})

	return r
}
