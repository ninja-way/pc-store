package repository

import (
	"context"
	"github.com/ninja-way/pc-store/internal/model"
)

// DB describes the repository business logic
type DB interface {
	GetComputers() ([]model.PC, error)
	GetComputerByID(int) (model.PC, error)
	AddComputer(model.PC) error
	UpdateComputer(int, model.PC) error
	DeleteComputer(int) error
	Close(context.Context) error
}
