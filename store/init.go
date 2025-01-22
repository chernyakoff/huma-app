package store

import (
	"database/sql"
	_ "embed"
	"errors"
	"log"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite" // SQLite driver for migration
	_ "github.com/golang-migrate/migrate/v4/source/file"     // Migration files
	"github.com/golang-migrate/migrate/v4/source/iofs"       // Migration files
	_ "modernc.org/sqlite"

	"huma-app/lib/config"
	"huma-app/store/migrations"
)

func Migrate() {
	d, err := iofs.New(migrations.FS, ".")
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithSourceInstance("embed://", d, "sqlite://"+config.Get().Storage.Path)
	if err != nil {
		slog.Error("cannot migrate db", "err", err)
		panic("cannot migrate db")
	}

	err = m.Up()
	if !errors.Is(err, migrate.ErrNoChange) {
		if err != nil {
			slog.Error("database migration failed",
				slog.Any("error", err),
				slog.String("database", "migrating: failure"))

			panic("database migration failed")
		}
		slog.Info("Database migrating: success")
	} else {
		slog.Info("Database migrating: no change required")
	}
	slog.Info("Database initialized")
}

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
	Migrate()

	return db
}
