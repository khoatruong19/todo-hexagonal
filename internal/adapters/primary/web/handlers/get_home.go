package handlers

import (
	"net/http"
	"todo-hexagonal/internal/adapters/primary/web/utils"
	pages "todo-hexagonal/internal/adapters/primary/web/views/pages/home"
)

type HomeHandler struct{}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	page := pages.Index()

	utils.RenderPage(w, r, "main", page, "Todoapp")
}
