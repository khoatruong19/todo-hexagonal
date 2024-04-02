package handlers

import (
	"net/http"
	layout "todo-hexagonal/internal/adapters/primary/web/views/layout/main"
	pages "todo-hexagonal/internal/adapters/primary/web/views/pages/home"
)

type HomeHandler struct{}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := pages.Index()
	err := layout.Layout(c, "My website").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
