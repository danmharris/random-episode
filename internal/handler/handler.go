package handler

import (
	"log/slog"
	"net/http"

	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/danmharris/random-episode/internal/ui"
)

type Handler struct {
	server     *http.Server
	tmdbClient *tmdb.Client
}

func NewHandler(tmdbClient *tmdb.Client) (*Handler, error) {
	ui.LoadTemplates()

	handler := &Handler{
		tmdbClient: tmdbClient,
	}

	router := http.NewServeMux()
	router.HandleFunc("GET /", handler.index)

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

	var shows []showData
	if q != "" {
		result, _ := h.tmdbClient.GetSearchTVShow(q, nil)

		for _, show := range result.Results {
			shows = append(shows, showData{
				Title: show.Name,
				ID:    show.ID,
			})
		}
	}

	w.WriteHeader(http.StatusOK)
	err := ui.RenderView("index.html.tmpl", w, struct {
		Title string
		Query string
		Shows []showData
	}{
		Title: "Home",
		Query: q,
		Shows: shows,
	})
	if err != nil {
		slog.Error("error rendering template", "error", err)
	}
}
