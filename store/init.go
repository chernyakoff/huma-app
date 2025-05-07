package store

import (
	"database/sql"
	"log/slog"

	_ "modernc.org/sqlite"

	"huma-app/lib/config"
)

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite", config.Get().Storage.Path)
	if err != nil {
		slog.Error("cannot open db connection", "err", err)
	}

	err = db.Ping()
	if err != nil {
		slog.Error("cannot ping db", "err", err)
	}

	slog.Info("Database connected: " + config.Get().Storage.Path)

	return db
}
