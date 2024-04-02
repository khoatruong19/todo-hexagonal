package ports

import (
	"todo-hexagonal/internal/core/domain"
)

type UserResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

type LoginUserResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
}

type UserRepository interface {
	CreateUser(email, username, password string) (*domain.User, error)
	GetUser(id string) (*domain.User, error)
	FindUserByEmail(email string) (*domain.User, error)
	FindUserByUsername(username string) (*domain.User, error)
	DeleteUser(id string) error
}

type UserService interface {
	CreateUser(email, username, password string) (*domain.User, error)
	GetUser(id string) (*domain.User, error)
	DeleteUser(id string) error
	LoginUser(username, password string) (LoginUserResponse, error)
}
