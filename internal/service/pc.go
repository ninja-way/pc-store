package service

import (
	"context"
	"github.com/ninja-way/pc-store/internal/models"
	"time"
)

type ComputerRepository interface {
	GetComputers(context.Context) ([]models.PC, error)
	GetComputerByID(context.Context, int) (models.PC, error)
	AddComputer(context.Context, models.PC) error
	UpdateComputer(context.Context, int, models.PC) error
	DeleteComputer(context.Context, int) error
}

type ComputersStore struct {
	repo ComputerRepository
}

func NewComputersStore(repo ComputerRepository) *ComputersStore {
	return &ComputersStore{
		repo: repo,
	}
}

func (c *ComputersStore) GetComputers(ctx context.Context) ([]models.PC, error) {
	return c.repo.GetComputers(ctx)
}

func (c *ComputersStore) GetComputerByID(ctx context.Context, i int) (models.PC, error) {
	return c.repo.GetComputerByID(ctx, i)
}

func (c *ComputersStore) AddComputer(ctx context.Context, pc models.PC) error {
	if pc.AddedAt.IsZero() {
		pc.AddedAt = time.Now()
	}

	return c.repo.AddComputer(ctx, pc)
}

func (c *ComputersStore) UpdateComputer(ctx context.Context, i int, pc models.PC) error {
	return c.repo.UpdateComputer(ctx, i, pc)
}

func (c *ComputersStore) DeleteComputer(ctx context.Context, i int) error {
	return c.repo.DeleteComputer(ctx, i)
}
