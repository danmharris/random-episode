package handler

import (
	"net/http"
)

func (h *handler) index(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title string
	}{
		Title: "Test",
	}

	h.views.RenderView(w, "index", data)
}
