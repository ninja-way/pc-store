package service

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ninja-way/pc-store/internal/models"
	"strconv"
	"time"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type UsersRepository interface {
	CreateUser(context.Context, models.User) error
	GetUser(context.Context, string, string) (models.User, error)
}

type Users struct {
	repo       UsersRepository
	hasher     PasswordHasher
	hmacSecret []byte
}

func NewUsers(repo UsersRepository, hasher PasswordHasher, secret []byte) *Users {
	return &Users{
		repo:       repo,
		hasher:     hasher,
		hmacSecret: secret,
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

func (u *Users) SignIn(ctx context.Context, inp models.SignIn) (string, error) {
	password, err := u.hasher.Hash(inp.Password)
	if err != nil {
		return "", err
	}

	user, err := u.repo.GetUser(ctx, inp.Email, password)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   strconv.Itoa(user.ID),
		IssuedAt:  &jwt.NumericDate{Time: time.Now()},
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Minute * 30)}, // TODO: move to config
	})

	return token.SignedString(u.hmacSecret)
}
