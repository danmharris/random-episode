package handler

import (
	"database/sql"
	"log/slog"
	"math/rand"
	"net/http"
	"strconv"
	"time"

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
	router.HandleFunc("POST /shows/{id}/episode", handler.createWatchedEpisode)

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

	result, _ := h.db.Exec("INSERT INTO shows (tmdb_id, title) VALUES ($1, $2)", tmdbID, showDetails.Name)
	id, _ := result.LastInsertId()
	http.Redirect(w, r, "/shows/"+strconv.Itoa(int(id)), http.StatusFound)
}

func (h *Handler) showShow(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	var show db.Show
	h.db.Get(&show, "SELECT * FROM shows WHERE id = $1", id)

	var watched []db.WatchedEpisode
	h.db.Select(&watched, "SELECT * FROM watched_episodes WHERE show_id = $1 ORDER BY timestamp DESC", id)

	w.WriteHeader(http.StatusOK)
	ui.RenderView("show.html.tmpl", w, struct {
		Title   string
		ID      int
		Watched []db.WatchedEpisode
	}{
		Title:   show.Title,
		ID:      show.ID,
		Watched: watched,
	})
}

func (h *Handler) showEpisode(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	var tmdbID int
	h.db.Get(&tmdbID, "SELECT tmdb_id FROM shows WHERE id = $1", id)

	showDetails, _ := h.tmdbClient.GetTVDetails(tmdbID, nil)

	var season int
	var episode int
	attempts := 0
	for {
		season = rand.Intn(showDetails.NumberOfSeasons)
		episode = rand.Intn(showDetails.Seasons[season].EpisodeCount)

		season++
		episode++

		var episodeId int
		err := h.db.Get(&episodeId, "SELECT id FROM watched_episodes WHERE show_id=$1 AND season=$2 AND episode=$3",
			r.PathValue("id"), season, episode)

		if err == sql.ErrNoRows {
			break
		}

		attempts++
		if attempts > 4 {
			slog.Warn("too many attempts to get episode that didn't clash, reusing latest")
			break
		}
	}

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

func (h *Handler) createWatchedEpisode(w http.ResponseWriter, r *http.Request) {
	showID, _ := strconv.Atoi(r.PathValue("id"))
	season, _ := strconv.Atoi(r.FormValue("season"))
	episode, _ := strconv.Atoi(r.FormValue("episode"))
	title := r.FormValue("title")

	timestamp := time.Now().Unix()

	h.db.Exec("INSERT INTO watched_episodes (show_id, season, episode, title, timestamp) VALUES ($1,$2,$3,$4,$5)",
		showID, season, episode, title, timestamp)

	http.Redirect(w, r, "/shows/"+r.PathValue("id"), http.StatusFound)
}
