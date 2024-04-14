package handlers

import (
	"net/http"
	"todo-hexagonal/internal/adapters/primary/web/httperror"
	"todo-hexagonal/internal/adapters/primary/web/utils"
	pages "todo-hexagonal/internal/adapters/primary/web/views/pages/home"
	"todo-hexagonal/internal/core/ports"
	"todo-hexagonal/internal/middleware"
)

type HomeHandler struct{}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(middleware.UserKey).(*ports.UserResponse)
	if !ok {
		httperror.ServerErrorResponse(w)
		return
	}

	page := pages.Index()

	utils.RenderPage(w, r, "main", page, "Todoapp")
}
