package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, statusCode int, body string, err error) {
	if err != nil {
		fmt.Println(err)
	}
	if statusCode > 499 {
		log.Printf("Responding with 5XX error: %s", body)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJson(w, statusCode, errorResponse{
		Error: body,
	})
}

func respondWithJson(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Unable to marhsal JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(statusCode)
	w.Write(data)
}
