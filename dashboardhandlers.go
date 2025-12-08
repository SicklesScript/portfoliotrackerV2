package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetDashboard(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userIDKey").(uuid.UUID)

	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unable to find user in context", nil)
		return
	}

	holdings, err := cfg.dbQueries.GetHoldings(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't load dashboard data", err)
		return
	}

	respondWithJson(w, http.StatusOK, holdings)
}
