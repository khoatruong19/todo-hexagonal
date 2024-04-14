package handlers

import (
	"net/http"
	"todo-hexagonal/internal/adapters/primary/web/httperror"
	"todo-hexagonal/internal/adapters/primary/web/utils"
	"todo-hexagonal/internal/constants"
	"todo-hexagonal/internal/core/services"
)

const (
	email_input_name           string = "email"
	username_input_name        string = "username"
	password_input_name        string = "password"
	confirmpassword_input_name string = "password"
)

type PostRegisterHandler struct {
	userService *services.UserService
}

type PostRegisterInput struct {
	Email           string `json:"email" validate:"required,min=5,max=20"`
	Username        string `json:"username" validate:"required"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirmPassword" validate:"required"`
}

func NewPostRegisterHandler(UserService *services.UserService) *PostRegisterHandler {
	return &PostRegisterHandler{
		userService: UserService,
	}
}

func (h *PostRegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue(email_input_name)
	username := r.FormValue(username_input_name)
	password := r.FormValue(password_input_name)
	confirmPassword := r.FormValue(confirmpassword_input_name)

	data := PostRegisterInput{
		Email:           email,
		Username:        username,
		Password:        password,
		ConfirmPassword: confirmPassword,
	}

	validationErrors := utils.ValidationInput(data)
	if validationErrors != nil {
		httperror.ValidationErrorResponse(w, validationErrors)
		return
	}

	_, err := h.userService.RegisterUser(email, username, password, confirmPassword)
	if err != nil {
		errMsg := err.Error()

		if errMsg == constants.ErrorEmailExisted {
			validationErrors := utils.NewValidationErrors(utils.NewValidationError(email_input_name, errMsg))
			httperror.ValidationErrorResponse(w, &validationErrors)
			return
		}

		if errMsg == constants.ErrorUsernameExisted {
			validationErrors := utils.NewValidationErrors(utils.NewValidationError(username_input_name, errMsg))
			httperror.ValidationErrorResponse(w, &validationErrors)
			return
		}

		if errMsg == constants.ErrorPasswordNotMatched {
			validationErrors := utils.NewValidationErrors(utils.NewValidationError(password_input_name, errMsg),
				utils.NewValidationError(confirmpassword_input_name, errMsg))
			httperror.ValidationErrorResponse(w, &validationErrors)
			return
		}

		httperror.ServerErrorResponse(w)
		return
	}

	w.Header().Set("HX-Redirect", "/login")
	w.WriteHeader(http.StatusCreated)
}
