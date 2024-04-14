package handlers

import (
	"net/http"
	"todo-hexagonal/internal/adapters/primary/web/httperror"

	"github.com/gorilla/sessions"
)

type PostLogoutHandler struct {
	sessionStore      *sessions.CookieStore
	sessionCookieName string
}

type NewPostLogoutParams struct {
	SessionStore      *sessions.CookieStore
	SessionCookieName string
}

func NewPostLogoutHandler(params NewPostLogoutParams) *PostLogoutHandler {
	return &PostLogoutHandler{
		sessionStore:      params.SessionStore,
		sessionCookieName: params.SessionCookieName,
	}
}

func (h *PostLogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session, err := h.sessionStore.Get(r, h.sessionCookieName)
	if err != nil {
		httperror.ServerErrorResponse(w)
		return
	}

	session.Options.MaxAge = -1

	err = session.Save(r, w)
	if err != nil {
		httperror.ServerErrorResponse(w)
		return
	}

	w.Header().Set("HX-Redirect", "/login")
	w.WriteHeader(http.StatusOK)
}
