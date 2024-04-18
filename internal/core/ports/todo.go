package ports

import (
	"todo-hexagonal/internal/core/domain"

	"github.com/google/uuid"
)

type TodoResponse struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Done   bool   `json:"done"`
	UserId string `json:"userId"`
}

type TodoRepository interface {
	CreateTodo(title string, userId uuid.UUID) (*domain.Todo, error)
	GetTodoById(id string) (*domain.Todo, error)
	GetTodosByUserId(userId uuid.UUID) (*[]domain.Todo, error)
	DeleteTodo(id string) error
}

type TodoService interface {
	CreateTodo(title string, userId uuid.UUID) (*TodoResponse, error)
	GetTodoById(id string) (*TodoResponse, error)
	GetTodosByUserId(userId uuid.UUID) (*[]TodoResponse, error)
	DeleteTodo(id string) error
}
