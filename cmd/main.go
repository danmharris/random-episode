package main

import (
	"os"

	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/danmharris/random-episode/internal/db"
	"github.com/danmharris/random-episode/internal/handler"
	"github.com/jmoiron/sqlx"
)

func main() {
	conn := sqlx.MustConnect("sqlite3", "file:random-episode.db")
	err := db.Migrate(conn.DB)
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
