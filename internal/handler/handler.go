package handler

import (
	"github.com/danmharris/random-episode/internal/data"
	"github.com/danmharris/random-episode/internal/view"
	"github.com/go-chi/chi/v5"
)

type handler struct {
	views view.Views
	tmdb  *data.TMDB
}

func Setup(tmdb *data.TMDB) chi.Router {
	router := chi.NewRouter()
	views := view.ParseViews()
	h := &handler{views, tmdb}

	router.Get("/", h.index)
	router.Get("/episode/{showID}", h.FindRandomEpisode)

	return router
}
