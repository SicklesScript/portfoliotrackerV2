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

// Function that serves login page
func (cfg *apiConfig) handlerLoginView(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/login.html")
}

// Function that validates user login attempt
func (cfg *apiConfig) handlerLoginAPI(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Attach token to json response for easier future authentication
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

/*
Function that validates user signup, logs users to database
User will be redirected to login page upon successful registration
*/
func (cfg *apiConfig) handlerSignupApi(w http.ResponseWriter, r *http.Request) {

	// Password and PasswordConfirmed are checked in signup.js
	type parameters struct {
		Email             string `json:"email"`
		Password          string `json:"password"`
		PasswordConfirmed string `json:"passwordConfirmed"`
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
		respondWithError(w, http.StatusConflict, "User already exists", err)
		return
	}

	respondWithJson(w, http.StatusOK, newUser)
}
