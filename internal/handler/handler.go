package handler

import (
	"github.com/danmharris/random-episode/internal/view"
	"github.com/go-chi/chi/v5"
)

type handler struct {
	views view.Views
}

func Setup() chi.Router {
	router := chi.NewRouter()
	views := view.ParseViews()
	h := &handler{views}

	router.Get("/", h.index)

	return router
}
