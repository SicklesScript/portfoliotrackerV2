package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/SicklesScript/portfoliotrackerV2/internal/auth"
	"github.com/SicklesScript/portfoliotrackerV2/internal/database"
)

func (cfg *apiConfig) handlerLoginView(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/login.html")
}

func (cfg *apiConfig) handlerLoginAPI(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"userEmail"`
		Password string `json:"userPassword"`
	}

	type response struct {
		ID        uuid.UUID `json:"id"`
		Email     string    `json:"email"`
		Token     string    `json:"token"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to decode json data", err)
		return
	}
	defer r.Body.Close()

	userData, err := cfg.dbQueries.GetUserByEmail(context.Background(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to retreive user data from db", err)
		return
	}

	match, err := auth.CheckHashedPassword(params.Password, userData.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to check password", err)
		return
	}
	if !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect password or email", err)
		return
	}
	token, err := auth.MakeJWT(userData.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create JWT token", err)
		return
	}

	respondWithJson(w, http.StatusOK, response{
		ID:        userData.ID,
		Email:     userData.Email,
		Token:     token,
		CreatedAt: userData.CreatedAt,
		UpdatedAt: userData.UpdatedAt,
	})
}

func (cfg *apiConfig) handlerSignupApi(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email             string `json:"userEmail"`
		Password          string `json:"userPass"`
		PasswordConfirmed string `json:"userPassConfirmed"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to decode json", err)
		return
	}

	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to hash password", err)
		return
	}

	newUser, err := cfg.dbQueries.CreateUser(context.Background(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPass,
	})
	if err != nil {
		respondWithError(w, http.StatusForbidden, "User already exists", err)
		return
	}

	respondWithJson(w, http.StatusOK, newUser)
}
