package handlers

import (
	"net/http"
	"todo-hexagonal/internal/adapters/primary/web/httperror"
	"todo-hexagonal/internal/adapters/primary/web/utils"
	pages "todo-hexagonal/internal/adapters/primary/web/views/pages/home"
	"todo-hexagonal/internal/core/ports"
	"todo-hexagonal/internal/core/services"
	"todo-hexagonal/internal/middleware"
)

type HomeHandler struct {
	todoService *services.TodoService
}

type NewHomeHandlerParams struct {
	TodoService *services.TodoService
}

func NewHomeHandler(params NewHomeHandlerParams) *HomeHandler {
	return &HomeHandler{
		todoService: params.TodoService,
	}
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserKey).(*ports.UserResponse)
	if !ok {
		httperror.ServerErrorResponse(w)
		return
	}

	todos, err := h.todoService.GetTodosByUserId(user.ID)
	if err != nil {
		httperror.ServerErrorResponse(w)
		return
	}

	page := pages.Index(user, *todos)

	utils.RenderPage(w, r, "main", page, "Todoapp")
}
