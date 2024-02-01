package repository

import (
	"context"
	"github.com/ninja-way/pc-store/internal/models"
)

// DB describes the repository business logic
type DB interface {
	// Users
	CreateUser(context.Context, models.User) error
	GetUser(context.Context, string, string) (models.User, error)

	// Tokens
	CreateToken(context.Context, models.RefreshSession) error
	GetToken(context.Context, string) (models.RefreshSession, error)

	// Computers
	GetComputers(context.Context) ([]models.PC, error)
	GetComputerByID(context.Context, int) (models.PC, error)
	AddComputer(context.Context, models.PC) (int, error)
	UpdateComputer(context.Context, int, models.PC) error
	DeleteComputer(context.Context, int) error

	Close(context.Context) error
}
