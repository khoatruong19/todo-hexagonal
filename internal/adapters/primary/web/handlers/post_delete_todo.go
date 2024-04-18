package handlers

import (
	"net/http"
	"todo-hexagonal/internal/adapters/primary/web/httperror"
	"todo-hexagonal/internal/core/ports"
	"todo-hexagonal/internal/core/services"
	"todo-hexagonal/internal/middleware"

	"github.com/go-chi/chi/v5"
)

type PostDeleteTodoHandler struct {
	todoService *services.TodoService
}

type NewPostDeleteTodoHandlerParams struct {
	TodoService *services.TodoService
}

func NewPostDeleteTodoHandler(params NewPostDeleteTodoHandlerParams) *PostDeleteTodoHandler {
	return &PostDeleteTodoHandler{
		todoService: params.TodoService,
	}
}

func (h *PostDeleteTodoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	todoId := chi.URLParam(r, "id")

	todo, err := h.todoService.GetTodoById(todoId)
	if err != nil || todo == nil {
		httperror.ServerErrorResponse(w)
		return
	}

	user, ok := r.Context().Value(middleware.UserKey).(*ports.UserResponse)
	if !ok {
		httperror.ServerErrorResponse(w)
		return
	}

	if todo.UserId != user.ID.String() {
		httperror.ServerErrorResponse(w)
		return
	}

	err = h.todoService.DeleteTodo(todoId)
	if err != nil {
		httperror.ServerErrorResponse(w)
		return
	}

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}
