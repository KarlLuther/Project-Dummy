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
	"github.com/joho/godotenv"
	"test.com/project/internal/models"
)

type application struct {
	cryptoKey []byte
	jwtSecret []byte
	logger    *slog.Logger
	secrets   *models.SecretModel
	users     *models.UserModel
}

func main() {
	addr := flag.String("addr", ":4000", "HTTPS network address")
	dsn := flag.String("dsn", "apiAdmin:ZyraVanya1337!@/secretstoragedb?parseTime=true", "MySQL data source name")
	flag.Parse()

	var encryptionKey string
	var jwtSecretKey string

	err := getDotEnv(&encryptionKey, &jwtSecretKey)
	if err != nil {
		log.Fatal(err)
	}

	if len(encryptionKey) != 32 {
		log.Fatal("Invalid encryption key size. Must be 32 bytes.")
	}

	if len(jwtSecretKey) != 512 {
		log.Fatal("Invalid JWT secret key size. Must be 512 bytes.")
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
		jwtSecret: []byte(jwtSecretKey),
		logger:    logger,
		secrets:   &models.SecretModel{DB: db},
		users:     &models.UserModel{DB: db},
	}

	router := app.routes()

	log.Printf("Starting HTTP server on %s...", *addr)
	logger.Info("starting HTTP server", "addr", *addr)

	err = http.ListenAndServe(*addr, router)
	if err != nil {
		logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
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

func getDotEnv(encryptionKey *string, jwtSecretKey *string) error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}

	*encryptionKey = os.Getenv("ENCRYPTION_KEY")
	if *encryptionKey == "" {
		return fmt.Errorf("ENCRYPTION_KEY not set in environment")
	}

	*jwtSecretKey = os.Getenv("JWT_SECRET")
	if *jwtSecretKey == "" {
		return fmt.Errorf("JWT_SECRET_KEY not set in environment")
	}

	return nil
}
