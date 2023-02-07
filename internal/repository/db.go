package repository

import (
	"context"
	"github.com/ninja-way/pc-store/internal/models"
)

// DB describes the repository business logic
type DB interface {
	GetComputers(context.Context) ([]models.PC, error)
	GetComputerByID(context.Context, int) (models.PC, error)
	AddComputer(context.Context, models.PC) (int, error)
	UpdateComputer(context.Context, int, models.PC) error
	DeleteComputer(context.Context, int) error
	Close(context.Context) error
}
