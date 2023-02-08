package service

import (
	"context"
	"github.com/ninja-way/pc-store/internal/models"
	"time"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type UsersRepository interface {
	CreateUser(context.Context, models.User) error
}

type Users struct {
	repo   UsersRepository
	hasher PasswordHasher
}

func NewUsers(repo UsersRepository, hasher PasswordHasher) *Users {
	return &Users{
		repo:   repo,
		hasher: hasher,
	}
}

func (u *Users) SignUp(ctx context.Context, inp models.SignUp) error {
	password, err := u.hasher.Hash(inp.Password)
	if err != nil {
		return err
	}

	user := models.User{
		Name:         inp.Name,
		Email:        inp.Email,
		Password:     password,
		RegisteredAt: time.Now(),
	}

	return u.repo.CreateUser(ctx, user)
}
