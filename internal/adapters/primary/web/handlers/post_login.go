package handlers

import (
	"net/http"
	"time"
	"todo-hexagonal/internal/adapters/primary/web/httperror"
	"todo-hexagonal/internal/adapters/primary/web/utils"
	"todo-hexagonal/internal/constants"
	"todo-hexagonal/internal/core/services"

	"github.com/gorilla/sessions"
)

type PostLoginHandler struct {
	userService       *services.UserService
	sessionStore      *sessions.CookieStore
	sessionCookieName string
}

type PostLoginInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type NewPostLoginParams struct {
	UserService       *services.UserService
	SessionStore      *sessions.CookieStore
	SessionCookieName string
}

func NewPostLoginHandler(params NewPostLoginParams) *PostLoginHandler {
	return &PostLoginHandler{
		userService:       params.UserService,
		sessionStore:      params.SessionStore,
		sessionCookieName: params.SessionCookieName,
	}
}

func (h *PostLoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue(username_input_name)
	password := r.FormValue(password_input_name)

	data := PostLoginInput{
		Username: username,
		Password: password,
	}

	validationErrors := utils.ValidationInput(data)
	if validationErrors != nil {
		httperror.ValidationErrorResponse(w, validationErrors)
		return
	}

	user, err := h.userService.LoginUser(username, password)
	if err != nil {
		errMsg := err.Error()

		if errMsg == constants.ErrorInvalidCredentials {
			validationErrors := utils.NewValidationErrors(utils.NewValidationError(username_input_name, errMsg),
				utils.NewValidationError(password_input_name, errMsg))
			httperror.ValidationErrorResponse(w, &validationErrors)
			return
		}

		httperror.ServerErrorResponse(w)
		return
	}

	session, _ := h.sessionStore.Get(r, h.sessionCookieName)
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60,
		HttpOnly: true,
	}

	session.Values = map[interface{}]interface{}{
		constants.AuthKey:     true,
		constants.UserIdKey:   user.ID,
		constants.TimezoneKey: time.Now().String(),
	}
	session.Save(r, w)

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}
