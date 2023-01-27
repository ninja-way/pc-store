package repository

import "github.com/ninja-way/pc-store/internal/model"

type Repository interface {
	GetComputers() ([]model.PC, error)
	GetComputerByID(int) (model.PC, error)
	AddComputer(model.PC) error
	UpdateComputer(int, model.PC) error
	DeleteComputer(int) error
}
