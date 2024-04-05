package handlers

import (
	"fmt"
	"net/http"
)

type PostLoginHandler struct{}

func NewPostLoginHandler() *PostLoginHandler {
	return &PostLoginHandler{}
}

func (h *PostLoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	fmt.Println(username)
	fmt.Println(password)

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}
