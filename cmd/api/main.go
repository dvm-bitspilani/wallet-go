package main

import (
	"dvm.wallet/harsh/cmd/api/config"
	"dvm.wallet/harsh/internal/database"
	"dvm.wallet/harsh/internal/version"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"log"
)

func main() {
	//logger := log.New(os.Stdout, "", log.LstdFlags|log.Llongfile)

	mainLogger, _ := zap.NewDevelopment()
	logger := mainLogger.Sugar()
	stdLogger := zap.NewStdLog(mainLogger)
	err := run(logger, stdLogger, mainLogger)
	if err != nil {
		logger.Fatal(err.Error())
	}
}

func run(logger *zap.SugaredLogger, stdLogger *log.Logger, mainLogger *zap.Logger) error {
	var cfg config.Config

	flag.StringVar(&cfg.BaseURL, "base-url", "http://localhost:4444", "base URL for the application")
	flag.IntVar(&cfg.HttpPort, "http-port", 4444, "port to listen on for HTTP requests")
	flag.StringVar(&cfg.Db.Dsn, "db-dsn", "host=127.0.0.1 port=5431 user=postgres dbname=wallet_db password=postgres sslmode=disable", "ent postgreSQL DSN")
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
		Config:     cfg,
		Client:     client,
		Logger:     logger,
		MainLogger: mainLogger,
		StdLogger:  stdLogger,
		Manager:    manager,
	}
	//ctx := context.Background()
	//pass, err := password.Hash("harsh")
	//u, _ := client.User.Query().Where(user.Username("vendorman")).Only(ctx)
	//app.Logger.Debugf(client.VendorSchema.Create().SetUser(u).SetClosed(false).SetName("vendy boi").SaveX(ctx).String())
	return serveHTTP(app)
}
