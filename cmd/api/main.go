package main

import (
	"dvm.wallet/harsh/cmd/api/config"
	"dvm.wallet/harsh/internal/database"
	"dvm.wallet/harsh/internal/version"
	"dvm.wallet/harsh/pkg/sse"
	"dvm.wallet/harsh/pkg/websocket"
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

	rdb, pubSub := sse.InitSSE()
	defer func() {
		// close entGo db connection
		err := client.Close()
		if err != nil {
			log.Fatal(err)
		}

		// Close the Redis connection
		err = rdb.Close()
		if err != nil {
			log.Fatal(err)
		}

		err = pubSub.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	manager := websocket.NewManager()

	app := &config.Application{
		Config:     cfg,
		Client:     client,
		Logger:     logger,
		MainLogger: mainLogger,
		StdLogger:  stdLogger,
		Manager:    manager,
		Rdb:        rdb,
		PubSub:     pubSub,
	}
	//app.Logger.Debugf(password.Hash("harsh"))
	//ctx := context.Background()
	//pass, err := password.Hash("harsh")
	//client.User.Create().
	//	SetUsername("harsh").
	//	SetEmail("harsh@gmail.com").
	//	SetUsername("f20211725").
	//	SetEmail("f20211725@pilani.bits-pilani.ac.in").
	//	SetPassword(pass).
	//	SetName("Harsh Singh").
	//	Save(ctx)

	//app.Logger.Debugf(client.VendorSchema.Create().SetUser(u).SetClosed(false).SetName("vendy boi").SaveX(ctx).String())
	return serveHTTP(app)
}
