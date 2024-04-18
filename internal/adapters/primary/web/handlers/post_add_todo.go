package handlers

import (
	"net/http"
	"todo-hexagonal/internal/adapters/primary/web/httperror"
	"todo-hexagonal/internal/adapters/primary/web/utils"
	"todo-hexagonal/internal/core/ports"
	"todo-hexagonal/internal/core/services"
	"todo-hexagonal/internal/middleware"
)

const (
	title_input_name string = "title"
)

type PostAddTodoHandler struct {
	todoService *services.TodoService
}

type NewPostAddTodoHandlerParams struct {
	TodoService *services.TodoService
}

type PostAddTodoInput struct {
	Title string `json:"title" validate:"max=40"`
}

func NewPostAddTodoHandler(params NewPostAddTodoHandlerParams) *PostAddTodoHandler {
	return &PostAddTodoHandler{
		todoService: params.TodoService,
	}
}

func (h *PostAddTodoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue(title_input_name)
	if title == "" {
		return
	}

	user, ok := r.Context().Value(middleware.UserKey).(*ports.UserResponse)
	if !ok {
		httperror.ServerErrorResponse(w)
		return
	}

	data := PostAddTodoInput{
		Title: title,
	}

	validationErrors := utils.ValidationInput(data)
	if validationErrors != nil {
		httperror.ValidationErrorResponse(w, validationErrors)
		return
	}

	_, err := h.todoService.CreateTodo(title, user.ID)
	if err != nil {
		httperror.ServerErrorResponse(w)
		return
	}

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusCreated)
}
