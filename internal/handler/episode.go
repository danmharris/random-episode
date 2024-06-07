package handler

import (
	"net/http"
	"strconv"

	"github.com/danmharris/random-episode/internal/episode"
	"github.com/go-chi/chi/v5"
)

func (h *handler) FindRandomEpisode(w http.ResponseWriter, r *http.Request) {
	type templateParams struct {
		Title   string
		Error   string
		Episode *episode.Episode
	}
	res := templateParams{}
	res.Title = "Episode"

	showIDParam := chi.URLParam(r, "showID")

	showID, err := strconv.Atoi(showIDParam)
	if err != nil {
		res.Error = "Invalid show ID"
		w.WriteHeader(http.StatusBadRequest)
		h.views.RenderView(w, "episode", res)
		return
	}

	episode, err := episode.FindRandomEpisode(h.tmdb, showID)
	if err != nil {
		res.Error = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		h.views.RenderView(w, "episode", res)
		return
	}

	res.Episode = episode
	h.views.RenderView(w, "episode", res)
}
