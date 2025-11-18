package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/SicklesScript/portfoliotrackerV2/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreatePortfolio(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userIDKey").(uuid.UUID)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unable to find user in context", nil)
		return
	}

	type parameters struct {
		PortfolioName string `json:"portfolioName"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to decode JSON", err)
		return
	}

	err = cfg.dbQueries.CreatePortfolio(context.Background(), database.CreatePortfolioParams{
		Name:   params.PortfolioName,
		UserID: userID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create portfolio in database", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
