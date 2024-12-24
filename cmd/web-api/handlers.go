package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (app *application) storeSecret(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Simulate parsing plaintext from the request
	plainText := r.URL.Query().Get("secret")
	if plainText == "" {
		http.Error(w, "Missing secret in query parameter", http.StatusBadRequest)
		return
	}


	secretID := app.secretNumber
	app.secretNumber++
	app.secrets[secretID] = plainText

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"status": "success",
		"secretID": fmt.Sprint(secretID),
	}
	json.NewEncoder(w).Encode(response)
}

func (app *application) getSecretByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	secretIDStr := vars["id"]

	secretIDInt, err := strconv.Atoi(secretIDStr)
	if err != nil {
		http.Error(w, "Invalid secret ID", http.StatusBadRequest)
		return
	}

	secret, ok := app.secrets[secretIDInt]
	if !ok {
		http.Error(w, "Secret not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"secretID": secretIDStr,
		"secret":   secret,
	}
	json.NewEncoder(w).Encode(response)
}

func (app *application) home(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("hello world"))
}