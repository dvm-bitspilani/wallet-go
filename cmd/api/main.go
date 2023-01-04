package wallet

import (
	"context"
	"dvm.wallet/harsh/cmd/api/config"
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

func run(logger *log.Logger) error {
	var cfg config.Config

	flag.StringVar(&cfg.BaseURL, "base-url", "http://localhost:4444", "base URL for the application")
	flag.IntVar(&cfg.HttpPort, "http-port", 4444, "port to listen on for HTTP requests")
	flag.StringVar(&cfg.Db.Dsn, "db-dsn", "host=127.0.0.1 port=5431 user=postgres dbname=wallet password=postgres sslmode=disable", "ent postgreSQL DSN")
	flag.BoolVar(&cfg.Db.Automigrate, "db-automigrate", true, "run migrations on startup")
	flag.StringVar(&cfg.Jwt.SecretKey, "jwt-secret-key", "rbztegymvi2bxjdh2tftkvd7b44z5akg", "secret key for JWT authentication")
	flag.BoolVar(&cfg.Version, "version", false, "display version and exit")

	flag.Parse()

	if cfg.Version {
		fmt.Printf("version: %s\n", version.Get())
		return nil
	}

	client, err := database.New(cfg.Db.Dsn, cfg.Db.Automigrate)
	if err != nil {
		return err
	}
	defer client.Close()

	app := &config.Application{
		Config: cfg,
		Client: client,
		Logger: logger,
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
	return serveHTTP(app)
}
