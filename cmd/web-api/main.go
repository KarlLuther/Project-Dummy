package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"test.com/project/internal/models"

	"github.com/joho/godotenv"
)

type application struct {
	cryptoKey []byte
	logger *slog.Logger
	secrets *models.SecretModel
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "apiAdmin:ZyraVanya1337!@/secretstoragedb?parseTime=true", "MySQL data source name")
	flag.Parse()

	var encryptionKey string 

	err := getDotEnv(&encryptionKey)
	if err != nil {
		log.Fatal(err)
	}

	if len(encryptionKey) != 32 {
    log.Fatal("Invalid encryption key size. Must be 32 bytes.")
	}	

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	app := &application{
		cryptoKey: []byte(encryptionKey),
		logger: logger, 
		secrets: &models.SecretModel{DB: db},
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

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil,err
	}

	//we don't actually create any connections with line 59, they are created
	//when needed, so we do a ping on the db in order to check that the connection 
	//will work. if there is an error we will close the connection pull and return the error
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func getDotEnv(encryptionKey *string) error{
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}

	*encryptionKey = os.Getenv("ENCRYPTION_KEY")
	if *encryptionKey == "" {
		return fmt.Errorf("ENCRYPTION_KEY not set in environment")
	}

	return nil
} 