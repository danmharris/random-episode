package main

import (
	"fmt"
	"net/http"

	"github.com/danmharris/random-episode/internal/config"
	"github.com/danmharris/random-episode/internal/data"
	"github.com/danmharris/random-episode/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.LoadConfigFromEnv()

	tmdb, err := data.NewTMDB(cfg.TMDBToken)
	if err != nil {
		panic(err)
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Mount("/", handler.Setup(tmdb))

	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), router)
}
