package main

import (
	"context"
	"net/http"

	"github.com/SicklesScript/portfoliotrackerV2/internal/auth"
)

/*
Middleware function that extracts and validates
JWT token from current user and passes the userID
along via r.WithContext()
*/
func (cfg *apiConfig) middlewareJWTAuthentication(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Auth error: No token", err)
			return
		}
		// Get userID via JWT validation
		claims, err := auth.ValidateJWT(token, cfg.jwtSecret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Auth error: Invalid token", err)
			return
		}
		// Create context with the userID value
		ctx := context.WithValue(r.Context(), "userIDKey", claims)

		// Return handler with userID attached
		handler(w, r.WithContext(ctx))
	})
}
