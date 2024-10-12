package db

import (
	"database/sql"
	"embed"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

type Show struct {
	ID     int
	TMDBID int `db:"tmdb_id"`
	Title  string
}

type WatchedEpisode struct {
	ID        int
	ShowID    int `db:"show_id"`
	Season    int
	Episode   int
	Title     string
	Timestamp int
}

//go:embed migrations/*.sql
var migrations embed.FS

func Migrate(conn *sql.DB, dbName string) error {
	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		return err
	}

	driver, err := pgx.WithInstance(conn, &pgx.Config{
		DatabaseName: dbName,
	})
	if err != nil {
		return err
	}

	migrate, err := migrate.NewWithInstance("embed", source, dbName, driver)
	if err != nil {
		return err
	}

	migrate.Up()
	return nil
}
