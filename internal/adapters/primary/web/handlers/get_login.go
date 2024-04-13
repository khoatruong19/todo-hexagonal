package handlers

import (
	"net/http"
	"todo-hexagonal/internal/adapters/primary/web/utils"
	pages "todo-hexagonal/internal/adapters/primary/web/views/pages/login"
)

type LoginHandler struct{}

func NewLoginHandler() *LoginHandler {
	return &LoginHandler{}
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	page := pages.Login()

	utils.RenderPage(w, r, "auth", page, "Login")
}
