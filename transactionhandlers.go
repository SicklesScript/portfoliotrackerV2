package main

import (
	"encoding/json"
	"net/http"

	"github.com/SicklesScript/portfoliotrackerV2/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateTransaction(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		PortfolioID   string `json:"portfolio_id"`
		StockName     string `json:"stock_name"`
		Ticker        string `json:"ticker"`
		Shares        string `json:"shares"`
		PricePerShare string `json:"price_per_share"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	portfolioUUID, err := uuid.Parse(params.PortfolioID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid portfolio ID format", err)
		return
	}

	err = cfg.dbQueries.CreateTransaction(r.Context(), database.CreateTransactionParams{
		StockName:     params.StockName,
		Ticker:        params.Ticker,
		Shares:        params.Shares,
		PricePerShare: params.PricePerShare,
		PortfolioID:   portfolioUUID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create transaction", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
