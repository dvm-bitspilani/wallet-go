package database

import (
	"context"
	"entgo.io/ent/dialect/sql/schema"
	"time"

	"dvm.wallet/harsh/ent"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const defaultTimeout = 3 * time.Second

type DB struct {
	*sqlx.DB
}

func New(dsn string, automigrate bool) (*ent.Client, error) {
	//db, err := sqlx.Connect("postgres", "postgres://"+dsn)
	//if err != nil {
	//	return nil, err
	//}
	//
	//db.SetMaxOpenConns(25)
	//db.SetMaxIdleConns(25)
	//db.SetConnMaxIdleTime(5 * time.Minute)
	//db.SetConnMaxLifetime(2 * time.Hour)
	//
	//if automigrate {
	//	iofsDriver, err := iofs.New(assets.EmbeddedFiles, "migrations")
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	migrator, err := migrate.NewWithSourceInstance("iofs", iofsDriver, "postgres://"+dsn)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	err = migrator.Up()
	//	switch {
	//	case errors.Is(err, migrate.ErrNoChange):
	//		break
	//	case err != nil:
	//		return nil, err
	//	}
	//}
	//
	//return &DB{db}, nil
	client, err := ent.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if automigrate {
		ctx := context.Background()
		err = client.Schema.Create(ctx, schema.WithAtlas(true)) // why the fuck is this deprecated?
		if err != nil {
			return nil, err
		}
	}

	return client, nil
}
