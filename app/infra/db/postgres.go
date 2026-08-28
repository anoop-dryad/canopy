package db

import (
	"github.com/anoop-dryad/canopy/app/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgresPool(cfg config.DB) *sqlx.DB {
	db, err := sqlx.Connect("postgres", cfg.DSN)
	if err != nil {
		panic("failed to connect to postgres: " + err.Error())
	}
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	return db
}
