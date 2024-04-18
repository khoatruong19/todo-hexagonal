package services

import (
	"todo-hexagonal/internal/core/ports"

	"github.com/google/uuid"
)

type TodoService struct {
	repo ports.TodoRepository
}

type NewTodoServiceParams struct {
	Repo ports.TodoRepository
}

func NewTodoService(params NewTodoServiceParams) *TodoService {
	return &TodoService{
		repo: params.Repo,
	}
}

func (t *TodoService) CreateTodo(title string, userId uuid.UUID) (*ports.TodoResponse, error) {
	todo, err := t.repo.CreateTodo(title, userId)
	if err != nil {
		return nil, err
	}

	return &ports.TodoResponse{
		ID:     todo.ID.String(),
		Title:  todo.Title,
		Done:   todo.Done,
		UserId: todo.UserID.String(),
	}, nil
}

func (t *TodoService) GetTodosByUserId(userId uuid.UUID) (*[]ports.TodoResponse, error) {
	todos, err := t.repo.GetTodosByUserId(userId)
	if err != nil {
		return nil, err
	}

	var result []ports.TodoResponse
	for _, todo := range *todos {
		result = append(result, ports.TodoResponse{
			ID:     todo.ID.String(),
			Title:  todo.Title,
			Done:   todo.Done,
			UserId: userId.String(),
		})
	}

	return &result, nil
}

func (t *TodoService) GetTodoById(id string) (*ports.TodoResponse, error) {
	todo, err := t.repo.GetTodoById(id)
	if err != nil {
		return nil, err
	}

	return &ports.TodoResponse{
		ID:     todo.ID.String(),
		Title:  todo.Title,
		Done:   todo.Done,
		UserId: todo.UserID.String(),
	}, nil
}

func (t *TodoService) DeleteTodo(id string) error {
	return t.repo.DeleteTodo(id)
}
