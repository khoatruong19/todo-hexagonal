package httperror

import (
	"net/http"
	"todo-hexagonal/internal/adapters/primary/web/response"
	"todo-hexagonal/internal/types"
)

func WriteError(w http.ResponseWriter, code int, error any) {
	err := response.WriteJSON(w, code, error, nil)
	if err != nil {
		w.WriteHeader(500)
	}
}

func RenderErrorResponse(w http.ResponseWriter) {
	message := "error rendering template"

	WriteError(w, http.StatusInternalServerError, message)
}

func ValidationErrorResponse(w http.ResponseWriter, validationErrors *[]types.ValidationError) {
	WriteError(w, http.StatusBadRequest, *validationErrors)
}

func UnauthorizedErrorResponse(w http.ResponseWriter) {
	message := "unauthorized"

	WriteError(w, http.StatusUnauthorized, message)
}

func ServerErrorResponse(w http.ResponseWriter) {
	message := "the server encountered a problem and could not process your request"

	WriteError(w, http.StatusInternalServerError, message)
}
