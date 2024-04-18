package services

import (
	"errors"
	"todo-hexagonal/internal/constants"
	"todo-hexagonal/internal/core/domain"
	"todo-hexagonal/internal/core/ports"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo ports.UserRepository
}

type NewUserServiceParams struct {
	Repo ports.UserRepository
}

func NewUserService(params NewUserServiceParams) *UserService {
	return &UserService{
		repo: params.Repo,
	}
}

func (u *UserService) CreateUser(email, username, password string) (*domain.User, error) {
	return u.repo.CreateUser(email, username, password)
}

func (u *UserService) GetUser(id string) (*ports.UserResponse, error) {
	user, err := u.repo.GetUser(id)
	if err != nil {
		return nil, err
	}

	return &ports.UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
		Avatar:   user.Avatar,
	}, nil
}

func (u *UserService) RegisterUser(email, username, password, confirmPassword string) (*ports.UserResponse, error) {
	existingEmail, _ := u.repo.FindUserByEmail(email)
	if existingEmail != nil {
		return nil, errors.New(constants.ErrorEmailExisted)
	}

	existingUsername, _ := u.repo.FindUserByUsername(username)
	if existingUsername != nil {
		return nil, errors.New(constants.ErrorUsernameExisted)
	}

	if password != confirmPassword {
		return nil, errors.New(constants.ErrorPasswordNotMatched)
	}

	createdUser, err := u.CreateUser(email, username, password)
	if err != nil {
		return nil, errors.New(constants.ErrorUserCreate)
	}

	return &ports.UserResponse{
		ID:       createdUser.ID,
		Email:    createdUser.Email,
		Username: createdUser.Username,
		Avatar:   createdUser.Avatar,
	}, nil
}

func (u *UserService) LoginUser(username, password string) (*ports.UserResponse, error) {
	user, err := u.repo.FindUserByUsername(username)
	if err != nil {
		return nil, errors.New(constants.ErrorInvalidCredentials)
	}

	err = u.verifyPassword(user.Password, password)
	if err != nil {
		return nil, err
	}

	return &ports.UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
		Avatar:   user.Avatar,
	}, nil
}

func (u *UserService) verifyPassword(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return errors.New(constants.ErrorInvalidCredentials)
	}
	return nil
}
