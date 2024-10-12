package main

import (
	"fmt"
	"os"

	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/danmharris/random-episode/internal/db"
	"github.com/danmharris/random-episode/internal/handler"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func main() {
	var (
		postgresUser string = os.Getenv("POSTGRES_USER")
		postgresPass string = os.Getenv("POSTGRES_PASS")
		postgresHost string = os.Getenv("POSTGRES_HOST")
		postgresPort string = os.Getenv("POSTGRES_PORT")
		postgresDB   string = os.Getenv("POSTGRES_DB")
	)

	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		postgresUser, postgresPass, postgresHost, postgresPort, postgresDB)

	conn := sqlx.MustConnect("pgx", connString)
	err := db.Migrate(conn.DB, postgresDB)
	if err != nil {
		panic(err)
	}

	tmdbClient, err := tmdb.InitV4(os.Getenv("TMDB_TOKEN"))
	if err != nil {
		panic(err)
	}

	handler, _ := handler.NewHandler(conn, tmdbClient)
	handler.Serve()
}
