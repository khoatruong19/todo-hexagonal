package repository

import (
	"errors"
	"fmt"
	"time"
	"todo-hexagonal/internal/core/domain"

	"golang.org/x/crypto/bcrypt"
)

func (u *DB) CreateUser(email, username, password string) (*domain.User, error) {
	user := &domain.User{}

	req := u.db.First(&user, "email = ? OR username = ?", email, username)
	if req.RowsAffected != 0 {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("password not hashed: %v", err)
	}

	user = &domain.User{
		Email:    email,
		Username: username,
		Password: string(hashedPassword),
	}

	req = u.db.Create(&user)
	if req.RowsAffected == 0 {
		return nil, fmt.Errorf("user not saved: %v", req.Error)
	}

	return user, nil
}

func (u *DB) GetUser(id string) (*domain.User, error) {
	user := &domain.User{}
	cachekey := user.ID.String()

	err := u.cache.Get(cachekey, &user)
	if err == nil {
		return user, nil
	}

	req := u.db.First(&user, "id = ? ", id)
	if req.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}

	err = u.cache.Set(cachekey, user, time.Minute*10)
	if err != nil {
		fmt.Printf("Error storing user in cache: %v", err)
	}

	return user, nil
}

func (u *DB) FindUserByEmail(email string) (*domain.User, error) {
	user := &domain.User{}

	req := u.db.First(&user, "email = ?", email)
	if req.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (u *DB) FindUserByUsername(username string) (*domain.User, error) {
	user := &domain.User{}

	req := u.db.First(&user, "username = ?", username)
	if req.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (u *DB) DeleteUser(id string) error {
	user := &domain.User{}

	req := u.db.Where("id = ?", id).Delete(&user)
	if req.RowsAffected == 0 {
		return errors.New("user not found")
	}

	err := u.cache.Delete(id)
	if err != nil {
		fmt.Printf("Error deleting user in cache: %v", err)
	}

	return nil
}
