package main

import (
	"github.com/gorilla/mux"
)

func (app *application) routes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/",app.home).Methods("GET")
	r.HandleFunc("/secrets/view/{id}", app.getSecretByID).Methods("GET")
	r.HandleFunc("/secrets/post", app.storeSecret).Methods("POST")
	r.HandleFunc("/register", app.registerNewUser).Methods("POST")
	r.HandleFunc("/login", app.login).Methods("POST")
	r.Use(app.authenticate)

	return r
} 