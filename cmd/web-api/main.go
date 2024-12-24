package main

import (
	"log"
	"net/http"
)

type application struct {
	secrets   map[int]string
	secretNumber int
}

func main() {
	app := &application{
		secrets: make(map[int]string),
		secretNumber: 1,
	}

	router := app.routes() 

	log.Println("Starting server on :4000...")
	err := http.ListenAndServe(":4000", router)
	if err != nil {
		log.Fatal(err)
	}
}
