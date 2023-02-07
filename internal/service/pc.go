package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/ninja-way/cache-ninja/pkg/cache"
	"github.com/ninja-way/pc-store/internal/config"
	"github.com/ninja-way/pc-store/internal/models"
	"time"
)

/******** Business logic layer *********/

const MaxPCPrice = 10000000

var (
	ErrPriceTooHigh  = errors.New("pc price too high")
	ErrFewComponents = errors.New("not all pc components listed")
)

// ComputerRepository is data layer entity
type ComputerRepository interface {
	GetComputers(context.Context) ([]models.PC, error)
	GetComputerByID(context.Context, int) (models.PC, error)
	AddComputer(context.Context, models.PC) (int, error)
	UpdateComputer(context.Context, int, models.PC) error
	DeleteComputer(context.Context, int) error
}

type ComputersStore struct {
	repo ComputerRepository

	cfg   *config.Config
	cache *cache.Cache
}

func NewComputersStore(c *cache.Cache, cfg *config.Config, repo ComputerRepository) *ComputersStore {
	return &ComputersStore{
		repo:  repo,
		cfg:   cfg,
		cache: c,
	}
}

func (c *ComputersStore) GetComputers(ctx context.Context) ([]models.PC, error) {
	return c.repo.GetComputers(ctx)
}

func (c *ComputersStore) GetComputerByID(ctx context.Context, id int) (models.PC, error) {
	pc, err := c.cache.Get(fmt.Sprintf("%d", id))
	if err == nil {
		return pc.(models.PC), nil
	}

	pc, err = c.repo.GetComputerByID(ctx, id)
	if err != nil {
		return models.PC{}, err
	}

	c.cache.Set(fmt.Sprintf("%d", id), pc, c.cfg.CacheTTL)
	return pc.(models.PC), nil
}

func (c *ComputersStore) AddComputer(ctx context.Context, pc models.PC) (int, error) {
	if pc.Price > MaxPCPrice {
		return 0, ErrPriceTooHigh
	}

	if pc.Name == "" || pc.CPU == "" || pc.RAM == 0 || pc.Price == 0 {
		return 0, ErrFewComponents
	}

	// if AddedAt value not specified, sets the current time
	if pc.AddedAt.IsZero() {
		pc.AddedAt = time.Now()
	}

	id, err := c.repo.AddComputer(ctx, pc)
	if err != nil {
		return 0, err
	}

	// add to cache
	pc.ID = id
	c.cache.Set(fmt.Sprintf("%d", id), pc, c.cfg.CacheTTL)
	return id, nil
}

func (c *ComputersStore) UpdateComputer(ctx context.Context, id int, newPC models.PC) error {
	if newPC.Price > MaxPCPrice {
		return ErrPriceTooHigh
	}

	// get old pc
	pc, err := c.repo.GetComputerByID(ctx, id)
	if err != nil {
		return err
	}

	// add to pc updates
	if newPC.Name != "" {
		pc.Name = newPC.Name
	}
	if newPC.CPU != "" {
		pc.CPU = newPC.CPU
	}
	if newPC.Videocard != "" {
		pc.Videocard = newPC.Videocard
	}
	if newPC.RAM != 0 {
		pc.RAM = newPC.RAM
	}
	if newPC.DataStorage != "" {
		pc.DataStorage = newPC.DataStorage
	}
	if newPC.Price != 0 {
		pc.Price = newPC.Price
	}

	if err = c.repo.UpdateComputer(ctx, id, pc); err != nil {
		return err
	}

	c.cache.Set(fmt.Sprintf("%d", id), pc, c.cfg.CacheTTL)
	return nil
}

func (c *ComputersStore) DeleteComputer(ctx context.Context, id int) error {
	return c.repo.DeleteComputer(ctx, id)
}
