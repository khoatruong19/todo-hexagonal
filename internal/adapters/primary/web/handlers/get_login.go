package handlers

import (
	"net/http"
	layout "todo-hexagonal/internal/adapters/primary/web/views/layout/auth"
	pages "todo-hexagonal/internal/adapters/primary/web/views/pages/login"
)

type LoginHandler struct{}

func NewLoginHandler() *LoginHandler {
	return &LoginHandler{}
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := pages.Login()
	err := layout.AuthLayout(c, "Login").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
