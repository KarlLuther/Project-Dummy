package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"test.com/project/internal/models"
)

//storeSecret is a handler for the POST /secrets/{name} endpoint
func (app *application) storeSecret(w http.ResponseWriter, r *http.Request) {
	//checking that the request method is appropriate
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	secretInstance,err := app.decodeJsonSecret(w,r)
	if err != nil {
		app.serverError(w,r,err)
		return
	}


	id, err := app.secrets.Insert(secretInstance.UserID, secretInstance.Name, secretInstance.SecretData, 7)
	if err != nil {
		app.serverError(w,r,err)
		return
	}

	app.writeJSONResponse(w, http.StatusOK, map[string]string{
    "status": "received",
    "id":     fmt.Sprint(id),
})
}

//getSecretByID is a handler for the GET /secrets/view/{id} endpoint
func (app *application) getSecretByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	secretIDStr := vars["id"]

	secretIDInt, err := strconv.Atoi(secretIDStr)
	if err != nil {
		app.logger.Warn("invalid secret ID")
		http.Error(w, "Invalid secret ID", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		app.clientError(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}


	secret, err := app.secrets.Get(secretIDInt,userID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.logger.Warn("secret not found", "secretID", secretIDInt)
			http.Error(w, "Secret not found", http.StatusNotFound)
		} else {
			app.logger.Error("failed to retrieve secret", "error", err)
			app.serverError(w, r, err)
		}
		return
	}

	plainText, err := app.decryptSecret(secret.SecretData)
	if err != nil {
		app.logger.Error("decryption failed", "secretID", secretIDStr, "error", err)
		app.serverError(w,r,err)
		return
	}

	app.logger.Info("secret retrieved succesfully", "secretID", secretIDStr)
	app.writeJSONResponse(w, http.StatusOK, map[string]string{
    "status": "success",
    "secretID": secretIDStr,
    "secret":   plainText,
})
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	//getting credentials out of the json body
	username, password, err := app.decodeJsonCredentials(w,r)
	if err != nil {
		app.serverError(w,r,err)
		return
	}

	//checking that the user actually exists within the db
	userID, err := app.users.Authenticate(username, password)
	if err != nil {
		app.clientError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	
	//generating a new JWT token string to send to the client
	tokenString, err := app.generateToken(w,r ,userID) 
	if err != nil {
		app.serverError(w,r,err)
		return
	}

	app.writeJSONResponse(w, http.StatusOK, map[string]string{
    "status": "success",
    "token":  tokenString,
})
}

func (app *application) registerNewUser(w http.ResponseWriter, r *http.Request) {
	//getting credentials out of the json body
	username, password, err := app.decodeJsonCredentials(w,r)
	if err != nil {
		app.serverError(w,r,err)
		return
	}

	//validating the password
	app.validatePassword(w,password)

	//hashPassword returns the bcrypt hash of a plaintext password
	//the second argument determines how strong the hash should be 
	//- the higher it is - the slower the hashing will be but more secure
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		app.serverError(w,r,err)
		return
	}

	userExists, err := app.users.UserExists(username)
	if err != nil {
		app.serverError(w,r,err)
		return
	}

	if userExists {
		app.clientError(w, "Username already in use", http.StatusConflict)
		return
	}

	//create a new user and insert it into the database
	id, err := app.users.Insert(username, hashedPassword)
	if err != nil {
		app.serverError(w,r,err)
		return
	}

	app.writeJSONResponse(w, http.StatusCreated, map[string]interface{}{
    "status": "success",
    "data":   map[string]int{"userID": id},
})
}

func (app *application) home(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("hello world"))
}

