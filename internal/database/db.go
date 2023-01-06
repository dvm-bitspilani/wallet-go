package database

import (
	"context"
	"dvm.wallet/harsh/ent"
	"entgo.io/ent/dialect/sql/schema"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/lib/pq"
)

const defaultTimeout = 3 * time.Second

func New(dsn string, automigrate bool) (*ent.Client, error) {
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
