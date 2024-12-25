package main

import (
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type application struct {
	secrets   map[int][]byte
	secretNumber int
	cryptoKey []byte
	logger *slog.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error while loading .env file")
	}

	encryptionKey := os.Getenv("ENCRYPTION_KEY")
	if encryptionKey == "" {
		log.Fatal("ENCRYPTION_KEY not found")
	}

	if len(encryptionKey) != 32 {
    log.Fatal("Invalid encryption key size. Must be 32 bytes.")
	}	

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	app := &application{
		secrets: make(map[int][]byte),
		secretNumber: 1,
		cryptoKey: []byte(encryptionKey),
		logger: logger, 
	}



	router := app.routes() 

	log.Println("Starting server on :4000...")
	logger.Info("stating server", "addr", *addr)

	err = http.ListenAndServe(*addr, router)
	if err != nil {
		logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}
