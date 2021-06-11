package schema

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func Connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func Migrate(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS shorturl (
		id serial primary key,
		created_at timestamp with time zone not null,
		deleted_at timestamp with time zone,
		url text not null unique,
		key text not null unique
	);`); err != nil {
		return err
	}
	return nil
}
