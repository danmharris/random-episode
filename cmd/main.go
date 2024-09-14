package main

import (
	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/danmharris/random-episode/internal/handler"
	"os"
)

func main() {
	tmdbClient, err := tmdb.InitV4(os.Getenv("TMDB_TOKEN"))
	if err != nil {
		panic(err)
	}

	handler, _ := handler.NewHandler(tmdbClient)
	handler.Serve()
}
