package model

import "database/sql"

type ModelClient struct {
	db *sql.DB
}

func NewClient(db *sql.DB) *ModelClient {
	return &ModelClient{
		db: db,
	}
}

func (c *ModelClient) ShortenMap() ShortenMap {
	return ShortenMap{
		db: c.db,
	}
}
