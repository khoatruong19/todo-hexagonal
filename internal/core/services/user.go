package services

import (
	"errors"
	"time"
	"todo-hexagonal/internal/config"
	"todo-hexagonal/internal/constants"
	"todo-hexagonal/internal/core/domain"
	"todo-hexagonal/internal/core/ports"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo ports.UserRepository
	cfg  *config.Config
}

type NewUserServiceParams struct {
	Repo ports.UserRepository
	Cfg  *config.Config
}

func NewUserService(params NewUserServiceParams) *UserService {
	return &UserService{
		repo: params.Repo,
		cfg:  params.Cfg,
	}
}

func (u *UserService) CreateUser(email, username, password string) (*domain.User, error) {
	return u.repo.CreateUser(email, username, password)
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
		ID:       createdUser.ID.String(),
		Email:    createdUser.Email,
		Username: createdUser.Username,
		Avatar:   createdUser.Avatar,
	}, nil
}

func (u *UserService) LoginUser(username, password string) (*ports.LoginUserResponse, error) {
	user, err := u.repo.FindUserByUsername(username)
	if err != nil {
		return nil, errors.New(constants.ErrorInvalidCredentials)
	}

	err = u.verifyPassword(user.Password, password)
	if err != nil {
		return nil, err
	}

	accessToken, err := u.generateAccessToken(user.ID.String())
	if err != nil {
		return nil, err
	}

	refreshToken, err := u.generateRefreshToken(user.ID.String())
	if err != nil {
		return nil, err
	}

	return &ports.LoginUserResponse{
		User: ports.UserResponse{
			ID:       user.ID.String(),
			Email:    user.Email,
			Username: user.Username,
			Avatar:   user.Avatar,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *UserService) verifyPassword(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return errors.New(constants.ErrorInvalidCredentials)
	}
	return nil
}

func (u *UserService) generateAccessToken(userID string) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "khoatruong19-access",
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour).UTC()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(u.cfg.JWTSecret))
}

func (u *UserService) generateRefreshToken(userID string) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "khoatruong19-refresh",
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour).UTC()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(u.cfg.JWTSecret))
}
