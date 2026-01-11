package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/yourusername/envsend/backend/internal/middleware"
	"github.com/yourusername/envsend/backend/internal/models"
	"github.com/yourusername/envsend/backend/internal/services"
)

// SecretsHandler handles secret-related HTTP requests.
type SecretsHandler struct {
	service *services.SecretService
}

// NewSecretsHandler creates a new secrets handler.
func NewSecretsHandler(service *services.SecretService) *SecretsHandler {
	return &SecretsHandler{service: service}
}

// CreateSecret handles POST /api/v1/secrets
func (h *SecretsHandler) CreateSecret(w http.ResponseWriter, r *http.Request) {
	var req models.CreateSecretRequest

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.EncryptedBlob == "" {
		http.Error(w, "encrypted_blob is required", http.StatusBadRequest)
		return
	}
	if req.EncryptionMetadata == nil {
		http.Error(w, "encryption_metadata is required", http.StatusBadRequest)
		return
	}
	if req.ExpiresIn == "" {
		http.Error(w, "expires_in is required", http.StatusBadRequest)
		return
	}
	if req.MaxViews < 1 {
		http.Error(w, "max_views must be at least 1", http.StatusBadRequest)
		return
	}

	// Get client IP
	ipAddress := getClientIP(r)

	// Create secret
	resp, err := h.service.CreateSecret(r.Context(), req, ipAddress)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// GetSecret handles GET /api/v1/secrets/:id
func (h *SecretsHandler) GetSecret(w http.ResponseWriter, r *http.Request) {
	secretID := chi.URLParam(r, "id")
	if secretID == "" {
		http.Error(w, "secret ID is required", http.StatusBadRequest)
		return
	}

	// Get client IP and user agent
	ipAddress := getClientIP(r)
	userAgent := r.UserAgent()

	// Get secret
	resp, err := h.service.GetSecret(r.Context(), secretID, ipAddress, userAgent)
	if err != nil {
		// Determine appropriate status code
		statusCode := http.StatusNotFound
		if err.Error() == "secret cannot be accessed (destroyed or expired)" {
			statusCode = http.StatusGone
		} else if err.Error() == "IP address mismatch" {
			statusCode = http.StatusForbidden
		}

		http.Error(w, err.Error(), statusCode)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// getClientIP extracts client IP from request.
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		for i := 0; i < len(xff); i++ {
			if xff[i] == ',' {
				return xff[:i]
			}
		}
		return xff
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}
