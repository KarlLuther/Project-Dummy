package main

import (
	"github.com/gorilla/mux"
)

func (app *application) routes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/",app.home).Methods("GET")
	r.HandleFunc("/secrets/view/{id}", app.getSecretByID).Methods("GET")
	r.HandleFunc("/secrets/{name}", app.storeSecret).Methods("POST")

	return r
} 