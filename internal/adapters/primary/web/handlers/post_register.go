package handlers

import (
	"net/http"
	pages "todo-hexagonal/internal/adapters/primary/web/views/pages/register"
	"todo-hexagonal/internal/core/services"
)

type PostRegisterHandler struct {
	userService *services.UserService
}

func NewPostRegisterHandler(UserService *services.UserService) *PostRegisterHandler {
	return &PostRegisterHandler{
		userService: UserService,
	}
}

func (h *PostRegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirmPassword")

	_, err := h.userService.RegisterUser(email, username, password, confirmPassword)
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

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("HX-Redirect", "/register")

	c := pages.RegisterSuccess()
	err = c.Render(r.Context(), w)
	if err != nil {
		http.Error(w, "error rendering template", http.StatusInternalServerError)
		return
	}
}
