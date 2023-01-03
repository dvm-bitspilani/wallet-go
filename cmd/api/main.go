package main

import (
	"context"
	"dvm.wallet/harsh/ent"
	"dvm.wallet/harsh/internal/password"
	"flag"
	"fmt"
	"log"
	"os"

	"dvm.wallet/harsh/internal/database"
	"dvm.wallet/harsh/internal/version"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Llongfile)

	err := run(logger)
	if err != nil {
		logger.Fatal(err)
	}
}

type config struct {
	baseURL  string
	httpPort int
	db       struct {
		dsn         string
		automigrate bool
	}
	jwt struct {
		secretKey string
	}
	version bool
}

type application struct {
	config config
	client *ent.Client
	logger *log.Logger
}

func run(logger *log.Logger) error {
	var cfg config

	flag.StringVar(&cfg.baseURL, "base-url", "http://localhost:4444", "base URL for the application")
	flag.IntVar(&cfg.httpPort, "http-port", 4444, "port to listen on for HTTP requests")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "host=127.0.0.1 port=5431 user=postgres dbname=wallet password=postgres sslmode=disable", "ent postgreSQL DSN")
	flag.BoolVar(&cfg.db.automigrate, "db-automigrate", true, "run migrations on startup")
	flag.StringVar(&cfg.jwt.secretKey, "jwt-secret-key", "rbztegymvi2bxjdh2tftkvd7b44z5akg", "secret key for JWT authentication")
	flag.BoolVar(&cfg.version, "version", false, "display version and exit")

	flag.Parse()

	if cfg.version {
		fmt.Printf("version: %s\n", version.Get())
		return nil
	}

	client, err := database.New(cfg.db.dsn, cfg.db.automigrate)
	if err != nil {
		return err
	}
	defer client.Close()

	app := &application{
		config: cfg,
		client: client,
		logger: logger,
	}
	//logger.Println(password.Hash("harsh"))
	ctx := context.Background()
	pass, err := password.Hash("harsh")
	logger.Println(client.User.Create().
		SetUsername("harsh").
		SetEmail("harsh@gmail.com").
		SetPassword(pass).
		SetName("Harsh Singh").
		Save(ctx))
	return app.serveHTTP()
}
