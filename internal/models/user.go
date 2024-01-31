package models

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"time"
)

// user input validator
var validate *validator.Validate

func init() {
	validate = validator.New()
}

var ErrUserNotFound = errors.New("user with these parameters not found")

type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	RegisteredAt time.Time `json:"registered_at"`
}

type SignUp struct {
	Name     string `json:"name" validate:"required,gte=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=6"`
}

func (s SignUp) Validate() error {
	return validate.Struct(s)
}

type SignIn struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=6"`
}

func (s SignIn) Validate() error {
	return validate.Struct(s)
}
