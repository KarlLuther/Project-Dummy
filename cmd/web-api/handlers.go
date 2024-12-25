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
		app.logger.Warn("invalid request method", "method", r.Method)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// Simulate parsing plaintext from the request
	plainText := r.URL.Query().Get("secret")
	if plainText == "" {
		app.logger.Warn("missing secret")
		app.clientError(w, http.StatusBadRequest)
		return
	}

	encryptedSecret, err := app.encryptSecret(plainText)
	if err != nil {
		app.logger.Error("encryption failed", "error", err)
		app.serverError(w,r,err)
		return
	}


	secretID := app.secretNumber
	app.secretNumber++
	app.secrets[secretID] = encryptedSecret

	app.logger.Info("secret stored succesfully", "secretID", secretID)
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
		app.logger.Warn("missing secret's ID")
		app.clientError(w, http.StatusBadRequest)
		return
	}

	encryptedSecret, ok := app.secrets[secretIDInt]
	if !ok {
		app.logger.Warn("secret not found", "secretID", secretIDStr)
		app.clientError(w, http.StatusNotFound)
		return
	}

	plainText, err := app.decryptSecret(encryptedSecret)
	if err != nil {
		app.logger.Error("decryption failed", "secretID", secretIDStr, "error", err)
		app.serverError(w,r,err)
	}

	app.logger.Info("secret retrieved succesfully", "secretID", secretIDStr)
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"secretID": secretIDStr,
		"secret":   plainText,
	}
	json.NewEncoder(w).Encode(response)
}

func (app *application) home(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("hello world"))
}