/**
 * @description
 * This file sets up the HTTP router for the Auth service using the Chi router.
 * It defines all the API routes, applies middleware like CORS and authentication,
 * and connects the routes to their respective handlers.
 *
 * @dependencies
 * - "net/http": For standard HTTP handling.
 * - "github.com/go-chi/chi/v5": The Chi router library.
 * - "github.com/go-chi/chi/v5/middleware": For standard Chi middleware.
 * - "github.com/go-chi/cors": For CORS middleware.
 */
package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// NewRouter creates and configures a new Chi router for the Auth service.
func NewRouter(handler *AuthHandler) http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Basic CORS configuration. This should be more restrictive in production.
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // Or specify your client's origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any major browsers
	}))

	// Health check endpoint - does not require authentication
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		// Apply the Clerk authentication middleware to this group.
		r.Use(ClerkAuth())

		// Define the onboarding route.
		r.Post("/onboarding", handler.OnboardingHandler)
	})

	return r
}