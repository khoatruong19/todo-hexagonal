package repository

import (
	"errors"
	"fmt"
	"todo-hexagonal/internal/core/domain"

	"github.com/google/uuid"
)

func (t *DB) CreateTodo(title string, userId uuid.UUID) (*domain.Todo, error) {
	todo := &domain.Todo{}

	todo = &domain.Todo{
		Title:  title,
		UserID: userId,
	}

	req := t.db.Create(&todo)
	if req.RowsAffected == 0 {
		return nil, fmt.Errorf("todo not saved: %v", req.Error)
	}

	return todo, nil
}

func (t *DB) GetTodosByUserId(userId uuid.UUID) (*[]domain.Todo, error) {
	todos := &[]domain.Todo{}

	_ = t.db.Find(&todos, "user_id = ?", userId.String())

	return todos, nil
}

func (t *DB) GetTodoById(id string) (*domain.Todo, error) {
	todo := &domain.Todo{}

	req := t.db.First(&todo, "id = ?", id)
	if req.RowsAffected == 0 {
		return nil, errors.New("todo not found")
	}

	return todo, nil
}

func (t *DB) DeleteTodo(id string) error {
	todo := &domain.Todo{}

	req := t.db.Where("id = ?", id).Delete(&todo)
	if req.RowsAffected == 0 {
		return errors.New("todo not found")
	}

	return nil
}
