package handlers

import (
	"net/http"
	layout "todo-hexagonal/internal/adapters/primary/web/views/layout/auth"
	pages "todo-hexagonal/internal/adapters/primary/web/views/pages/register"
)

type RegisterHandler struct{}

func NewRegisterHandler() *RegisterHandler {
	return &RegisterHandler{}
}

func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := pages.Register()
	err := layout.AuthLayout(c, "Register").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
