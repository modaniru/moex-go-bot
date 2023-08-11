package storage

import (
	"database/sql"

	gen "github.com/modaniru/moex-telegram-bot/internal/storage/gen"
)

type Storager interface {
	gen.Querier
}

type Storage struct {
	gen.Querier
	db *sql.DB
}

func NewStorage(source *sql.DB, q gen.Querier) *Storage {
	return &Storage{
		Querier: q,
		db:      source,
	}
}
