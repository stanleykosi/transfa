/**
 * @description
 * This file contains the HTTP handlers for the Auth service. Handlers are responsible
 * for parsing incoming requests, calling the appropriate application service method,
 * and writing the HTTP response.
 *
 * @dependencies
 * - "encoding/json": For JSON serialization and deserialization.
 * - "log": For logging.
 * - "net/http": For standard HTTP handling.
 * - "github.com/clerk/clerk-sdk-go/v2/jwt": To access session claims.
 * - "transfa/services/auth/internal/app": Imports the application service layer.
 * - "transfa/services/auth/internal/domain": Imports the data models/DTOs.
 */
package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"transfa/services/auth/internal/app"
	"transfa/services/auth/internal/domain"
)

// AuthHandler holds dependencies for the authentication-related HTTP handlers.
type AuthHandler struct {
	service *app.Service
}

// NewAuthHandler creates a new handler with the given application service.
func NewAuthHandler(service *app.Service) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

// OnboardingHandler handles the `POST /onboarding` request.
func (h *AuthHandler) OnboardingHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Get claims from context (set by middleware).
	claims, ok := r.Context().Value(sessionClaimsKey).(*jwt.Claims)
	if !ok || claims == nil {
		http.Error(w, "Unauthorized: Could not retrieve claims", http.StatusUnauthorized)
		return
	}

	// The Clerk User ID is in the 'Subject' field of the claims.
	clerkID := claims.Subject

	// 2. Decode the request body.
	var req domain.OnboardingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request: Invalid JSON body", http.StatusBadRequest)
		return
	}

	// 3. Basic validation (can be expanded).
	if req.Username == "" || req.AccountType == "" {
		http.Error(w, "Bad Request: username and account_type are required", http.StatusBadRequest)
		return
	}
	if req.AccountType != "personal" && req.AccountType != "merchant" {
		http.Error(w, "Bad Request: account_type must be 'personal' or 'merchant'", http.StatusBadRequest)
		return
	}

	// 4. Call the application service.
	user, err := h.service.OnboardUser(r.Context(), clerkID, req)
	if err != nil {
		// A more sophisticated error handling could map service errors to HTTP status codes.
		// For example, a "username taken" error could map to 409 Conflict.
		log.Printf("Onboarding failed for clerk_id %s: %v", clerkID, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 5. Write the success response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}