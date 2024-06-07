package handler

import (
	"github.com/danmharris/random-episode/internal/episode"
	"github.com/danmharris/random-episode/internal/view"
	"github.com/go-chi/chi/v5"
)

type handler struct {
	views    view.Views
	showData episode.EpisodeData
}

func Setup() chi.Router {
	router := chi.NewRouter()
	views := view.ParseViews()
	ed := episodeData{}
	h := &handler{views, &ed}

	router.Get("/", h.index)
	router.Get("/episode/{showID}", h.FindRandomEpisode)

	return router
}

// Temporary stub struct & implementation
type episodeData struct{}

func (s *episodeData) ShowDetails(id int) (*episode.Show, error) {
	episodeCounts := []int{25, 25, 25, 25, 25}
	return &episode.Show{
		Title:        "Placeholder Show",
		EpisodeCount: episodeCounts,
	}, nil
}

func (s *episodeData) EpisodeDetails(sh *episode.Show, pos episode.EpisodePosition) (*episode.Episode, error) {
	return &episode.Episode{
		EpisodePosition: pos,
		Title:           "An Episode",
		Show:            sh,
	}, nil
}
