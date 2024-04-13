package handlers

import (
	"net/http"
	"todo-hexagonal/internal/adapters/primary/web/utils"
	pages "todo-hexagonal/internal/adapters/primary/web/views/pages/register"
)

type RegisterHandler struct{}

func NewRegisterHandler() *RegisterHandler {
	return &RegisterHandler{}
}

func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	page := pages.Register()

	utils.RenderPage(w, r, "auth", page, "Register")
}
