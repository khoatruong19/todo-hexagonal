package handlers

import (
	"net/http"
	pages "todo-hexagonal/internal/adapters/primary/web/views/pages/register"
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

		c := pages.RegisterError(err.Error())
		err = c.Render(r.Context(), w)
		if err != nil {
			http.Error(w, "error rendering template", http.StatusInternalServerError)
			return
		}

		return
	}

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}
