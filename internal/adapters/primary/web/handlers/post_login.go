package handlers

import (
	"net/http"
	"todo-hexagonal/internal/core/services"
)

type PostLoginHandler struct {
	userService *services.UserService
}

func NewPostLoginHandler(UserService *services.UserService) *PostLoginHandler {
	return &PostLoginHandler{
		userService: UserService,
	}
}

func (h *PostLoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	_, err := h.userService.LoginUser(username, password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}
