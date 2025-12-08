package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SicklesScript/portfoliotrackerV2/internal/database"
	"github.com/google/uuid"
)

/*
Handler that alows users to create a new portfolio
*/
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

func (cfg *apiConfig) handlerGetPortfolios(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userIDKey").(uuid.UUID)

	portfolios, err := cfg.dbQueries.GetPortfoliosForUser(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve portfolios", err)
		return
	}

	respondWithJson(w, http.StatusOK, portfolios)
}

type FMPKeyMetrics struct {
	Symbol                  string  `json:"symbol"`
	ReturnOnEquity          float64 `json:"returnOnEquity"`
	ReturnOnInvestedCapital float64 `json:"returnOnInvestedCapital"`
	CurrentRatio            float64 `json:"currentRatio"`
}

/*
Retreives a ticker from dashboard.js and queries the
fmp api to send back information for each holding
*/
func (cfg *apiConfig) handlerGetKeyMetrics(w http.ResponseWriter, r *http.Request) {
	ticker := r.URL.Query().Get("ticker")
	if ticker == "" {
		respondWithError(w, http.StatusBadRequest, "Ticker is required", nil)
		return
	}

	apiKey := cfg.fmpApiKey

	url := fmt.Sprintf("https://financialmodelingprep.com/stable/key-metrics?symbol=%s&apikey=%s", ticker, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to contact FMP API", err)
		return
	}
	defer resp.Body.Close()

	var metrics []FMPKeyMetrics
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&metrics); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to parse FMP data", err)
		return
	}

	if len(metrics) == 0 {
		respondWithError(w, http.StatusNotFound, "No metrics found for ticker", err)
		return
	}

	respondWithJson(w, http.StatusOK, metrics[0])
}
