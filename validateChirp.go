package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func handlerValidation(w http.ResponseWriter, req *http.Request) {

	type ChirpRequest struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(req.Body)
	chirpReq := ChirpRequest{}
	err := decoder.Decode(&chirpReq)
	if err != nil {
		log.Printf("Error decoding chirp request: %s", err)
		w.WriteHeader(500)
		return
	}

	maxChirpLength := 140
	successCode := 200
	errorCode := 400
	wrongMethodCode := 405

	badWords := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}

	if req.Method != "POST" {
		log.Printf("Method must be POST: %s", req.Method)
		respondWithError(w, wrongMethodCode, "Method not allowed")
	}

	if len(chirpReq.Body) > maxChirpLength {
		respondWithError(w, errorCode, "Chirp too long")
	} else {
		words := strings.Split(chirpReq.Body, " ")
		for _, word := range badWords {
			for i, chirpWord := range words {
				if strings.ToLower(chirpWord) == word {
					words[i] = "****"
				}
			}
		}
		cleanedChirp := strings.Join(words, " ")
		respondWithJSON(w, successCode, map[string]string{"cleaned_body": cleanedChirp})
		return
	}

}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type ErrorResponse struct {
		Error string `json:"error"`
	}
	resp := ErrorResponse{
		Error: msg,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error marshaling data: %s", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(data))
	return

}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {

	data, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, 400, "Error handling json")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(data))

	return
}
