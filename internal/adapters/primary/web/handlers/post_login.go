package handlers

import (
	"net/http"
	"todo-hexagonal/internal/adapters/primary/web/httperror"
	"todo-hexagonal/internal/adapters/primary/web/utils"
	"todo-hexagonal/internal/constants"
	"todo-hexagonal/internal/core/services"
)

type PostLoginHandler struct {
	userService *services.UserService
}

type PostLoginInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func NewPostLoginHandler(UserService *services.UserService) *PostLoginHandler {
	return &PostLoginHandler{
		userService: UserService,
	}
}

func (h *PostLoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	data := PostLoginInput{
		Username: username,
		Password: password,
	}

	validationErrors := utils.ValidationInput(data)
	if validationErrors != nil {
		httperror.ValidationErrorResponse(w, validationErrors)
		return
	}

	_, err := h.userService.LoginUser(username, password)
	if err != nil {
		errMsg := err.Error()

		if errMsg == constants.ErrorInvalidCredentials {
			validationErrors := utils.NewValidationErrors(utils.NewValidationError("Username", errMsg), utils.NewValidationError("Password", errMsg))
			httperror.ValidationErrorResponse(w, &validationErrors)
			return
		}

		httperror.ServerErrorResponse(w)
		return
	}

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}
