package model

import (
	"context"
	"database/sql"
	"time"
)

type ShortenMap struct {
	db        *sql.DB
	ID        int64     `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
	DeletedAt time.Time `json:"deletedAt"`
	Url       string    `json:"url"`
	Key       string    `json:"key"`
}

func (s ShortenMap) Create(ctx context.Context, url, key string) (*ShortenMap, error) {
	insertStmt, err := s.db.PrepareContext(ctx, "insert into shorturl(created_at, url, key) values($1, $2, $3)")
	if err != nil {
		return nil, err
	}

	createdAt := time.Now()
	if _, err := insertStmt.ExecContext(ctx, createdAt, url, key); err != nil {
		return nil, err
	}

	return &ShortenMap{db: s.db, CreatedAt: createdAt, Url: url, Key: key}, nil
}

func (s ShortenMap) FindOne(ctx context.Context, key string) (*ShortenMap, error) {
	row := s.db.QueryRowContext(ctx, "select id, created_at, url from shorturl where key=$1 and deleted_at is null", key)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var (
		id        int64
		url       string
		createdAt time.Time
	)
	if err := row.Scan(&id, &createdAt, &url); err != nil {
		return nil, err
	}
	return &ShortenMap{
		ID:        id,
		CreatedAt: createdAt,
		Url:       url,
		Key:       key,
	}, nil
}

func (s *ShortenMap) Delete(ctx context.Context) error {
	deletedAt := time.Now()
	if _, err := s.db.ExecContext(ctx, "update shorturl set deleted_at=$1 where key=$2", deletedAt, s.Key); err != nil {
		return err
	}
	return nil
}
