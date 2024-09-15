package handler

import (
	"log/slog"
	"math/rand"
	"net/http"
	"strconv"

	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/danmharris/random-episode/internal/db"
	"github.com/danmharris/random-episode/internal/ui"
	"github.com/jmoiron/sqlx"
)

type Handler struct {
	server     *http.Server
	tmdbClient *tmdb.Client
	db         *sqlx.DB
}

func NewHandler(conn *sqlx.DB, tmdbClient *tmdb.Client) (*Handler, error) {
	ui.LoadTemplates()

	handler := &Handler{
		tmdbClient: tmdbClient,
		db:         conn,
	}

	router := http.NewServeMux()
	router.HandleFunc("GET /", handler.index)
	router.HandleFunc("POST /shows", handler.createShow)
	router.HandleFunc("GET /shows/{id}", handler.showShow)
	router.HandleFunc("GET /shows/{id}/episode", handler.showEpisode)

	handler.server = &http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	return handler, nil
}

func (h *Handler) Serve() {
	slog.Info("starting server", "address", h.server.Addr)
	h.server.ListenAndServe()
}

func (h *Handler) index(w http.ResponseWriter, r *http.Request) {
	type showData struct {
		Title string
		ID    int64
	}
	q := r.FormValue("q")

	var searchShows []showData
	if q != "" {
		result, _ := h.tmdbClient.GetSearchTVShow(q, nil)

		for _, show := range result.Results {
			searchShows = append(searchShows, showData{
				Title: show.Name,
				ID:    show.ID,
			})
		}
	}

	var shows []db.Show
	h.db.Select(&shows, "SELECT * FROM shows")

	w.WriteHeader(http.StatusOK)
	err := ui.RenderView("index.html.tmpl", w, struct {
		Title         string
		Query         string
		SearchResults []showData
		Shows         []db.Show
	}{
		Title:         "Home",
		Query:         q,
		SearchResults: searchShows,
		Shows:         shows,
	})
	if err != nil {
		slog.Error("error rendering template", "error", err)
	}
}

func (h *Handler) createShow(w http.ResponseWriter, r *http.Request) {
	tmdbID, _ := strconv.Atoi(r.FormValue("id"))

	showDetails, err := h.tmdbClient.GetTVDetails(tmdbID, nil)
	if err != nil {
		slog.Error("error getting show details", "error", err)
	}

	result, _ := h.db.Exec("INSERT INTO shows (tmdb_id, title) VALUES (?, ?)", tmdbID, showDetails.Name)
	id, _ := result.LastInsertId()
	http.Redirect(w, r, "/shows/"+strconv.Itoa(int(id)), http.StatusFound)
}

func (h *Handler) showShow(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	var show db.Show
	h.db.Get(&show, "SELECT * FROM shows WHERE id = ?", id)

	w.WriteHeader(http.StatusOK)
	ui.RenderView("show.html.tmpl", w, struct {
		Title string
		ID    int
	}{
		Title: show.Title,
		ID:    show.ID,
	})
}

func (h *Handler) showEpisode(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	var tmdbID int
	h.db.Get(&tmdbID, "SELECT tmdb_id FROM shows WHERE id = ?", id)

	showDetails, _ := h.tmdbClient.GetTVDetails(tmdbID, nil)

	season := rand.Intn(showDetails.NumberOfSeasons)
	episode := rand.Intn(showDetails.Seasons[season].EpisodeCount)

	season++
	episode++

	episodeDetails, _ := h.tmdbClient.GetTVEpisodeDetails(tmdbID, season, episode, nil)

	w.WriteHeader(http.StatusOK)
	ui.RenderView("episode.html.tmpl", w, struct {
		Title   string
		Season  int
		Episode int
		Path    string
	}{
		Title:   episodeDetails.Name,
		Season:  season,
		Episode: episode,
		Path:    r.URL.Path,
	})
}
