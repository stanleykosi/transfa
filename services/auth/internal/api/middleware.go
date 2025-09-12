/**
 * @description
 * This file contains the authentication middleware for the Auth service's API.
 * It uses the Clerk Go SDK to validate JWTs from incoming requests, ensuring
 * that only authenticated users can access protected endpoints.
 *
 * @dependencies
 * - "context": To manage request-scoped values like session claims.
 * - "net/http": For standard HTTP handling.
 * - "strings": For string manipulation.
 * - "github.com/clerk/clerk-sdk-go/v2": The official Clerk SDK.
 * - "github.com/clerk/clerk-sdk-go/v2/jwt": For JWT verification.
 */
package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
)

// claimsContextKey is a custom type to use as a key for storing claims in the request context.
type claimsContextKey string

const sessionClaimsKey claimsContextKey = "session_claims"

// ClerkAuth is a middleware that validates the Clerk session token from the Authorization header.
func ClerkAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the session token from the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Unauthorized: Missing Authorization Header", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")

			// Verify the token
			claims, err := jwt.Verify(r.Context(), &jwt.VerifyParams{
				Token: token,
			})
			if err != nil {
				http.Error(w, "Unauthorized: Invalid Token", http.StatusUnauthorized)
				return
			}

			// Add the claims to the request context
			ctx := context.WithValue(r.Context(), sessionClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}