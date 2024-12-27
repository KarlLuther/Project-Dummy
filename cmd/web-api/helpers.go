package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"test.com/project/internal/models"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		url = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "url", url)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, errorText string,status int) {
	http.Error(w, errorText, status)
} 

func (app *application) writeJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (app *application) decodeJsonSecret(w http.ResponseWriter, r *http.Request) (*models.Secret, error){
	/*
	we need this intermediate stage where we check whether data provided in
	the json can be converted into data secretInstance expects, 
	since go doesn't automatically convert variables unless it's
	explicitaly specified 
	*/

	var jsonBody struct {
		UserID     string `json:"UserID"`
		Name       string `json:"Name"`
		SecretData string `json:"SecretData"`
	}

	err := json.NewDecoder(r.Body).Decode(&jsonBody)
	if err != nil {
		app.serverError(w, r, err)
		return nil, err
	}
			// Check for missing SecretData
			if jsonBody.SecretData == "" || jsonBody.Name == ""{
				app.clientError(w, "Missing SecretData or Name", http.StatusBadRequest)
					return nil, fmt.Errorf("missing Name or Secret")
			}

			encryptedData, err := app.encryptSecret(jsonBody.SecretData)
			if err != nil {
				app.serverError(w, r, err)
				return nil, err
			}
	
			// Create Secret instance and assign values
			userID, err := strconv.Atoi(jsonBody.UserID)
			if err != nil {
				app.serverError(w, r, err)
				return nil, err
			}

			secretInstance := models.Secret{
					UserID:     userID,
					Name:       jsonBody.Name,
					SecretData: encryptedData,
			}

	return &secretInstance, nil
}

func (app *application) decodeJsonCredentials(w http.ResponseWriter, r *http.Request) (string, string, error){
	//defining a credentionals json struct to store requests user's credentials
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	//decoding the json body to populate credentials struct with the payload values
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		app.clientError(w, "Invalid credentials", http.StatusBadRequest)
		return "", "", err
	}

	// Checking if the username or password are missing
	if credentials.Username == "" || credentials.Password == "" {
		app.clientError(w, "Missing username or password", http.StatusBadRequest)
		return "", "", fmt.Errorf("missing credentials")
	}

	return credentials.Username, credentials.Password, nil
}


func (app *application) generateToken(w http.ResponseWriter, r *http.Request, userID int) (string, error) {
	//creating a new jwt token. It's hashed usig HS256 algorithm. 
	//claims are credentials attached to the token. Expiration is set to 24 hours
	//the token still needs to be signed using a secret, otherwise it won't be valid and can be tempered with
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	//sealing the token, making sure that it won't be tampered with by unintended users
	tokenString, err := token.SignedString(app.jwtSecret)
	if err != nil {
		app.serverError(w, r, err)
		return  "", err
	}

	return tokenString, nil
}

func (app *application) validatePassword(w http.ResponseWriter, password string) {
	//checking if the password is of sufficient length 
	hasMinLength := len(password) >= 8
	//these are password criterias that are false by default, and will be set to true if they are met
	hasLetter := false
	hasNumber := false
	hasSpecial := false
	hasUpper := false

	//checking every character of the password. if at least one 
	//character meets on the criterias it will be set to true
	for _, char := range password {
		switch {
		case 'a' <= char && char <= 'z':
			hasLetter = true
		case 'A' <= char && char <= 'Z':
			hasLetter = true
			hasUpper = true
		case '0' <= char && char <= '9':
			hasNumber = true
		case char >= 33 && char <= 47 || char >= 58 && char <= 64 || char >= 91 && char <= 96 || char >= 123 && char <= 126:
			hasSpecial = true
		}
	}

	//if at least one of the password criterias is not met the request will be rejected
	if !hasMinLength || !hasLetter || !hasNumber || !hasSpecial || !hasUpper {
		app.clientError(w, "Password must be at least 8 characters long and include a letter, a number, a special character, and an uppercase letter", http.StatusBadRequest)
		return
	}
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.clientError(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := authHeader[len("Bearer "):]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return app.jwtSecret, nil
		})

		if err != nil || !token.Valid {
			app.clientError(w, "Invalid token", http.StatusUnauthorized)
			return
		}	

		next.ServeHTTP(w, r)
	})
}
