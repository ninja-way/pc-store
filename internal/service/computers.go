package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/ninja-way/cache-ninja/pkg/cache"
	audit "github.com/ninja-way/grpc-audit-log/pkg/models"
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

type ComputersStore struct {
	repo     ComputerRepository
	auditLog AuditClient

	cfg   *config.Config
	cache *cache.Cache
}

func NewComputersStore(c *cache.Cache, cfg *config.Config, repo ComputerRepository, auditClient AuditClient) *ComputersStore {
	return &ComputersStore{
		repo:     repo,
		auditLog: auditClient,
		cfg:      cfg,
		cache:    c,
	}
}

func (c *ComputersStore) GetComputers(ctx context.Context) ([]models.PC, error) {
	comps, err := c.repo.GetComputers(ctx)
	if err != nil {
		return nil, err
	}

	if err = c.auditLog.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ACTION_GET,
		Entity:    audit.ENTITY_COMPUTER,
		EntityID:  0,
		Timestamp: time.Now(),
	}); err != nil {
		config.LogError("getComputers", err)
	}

	return comps, nil
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
	gotPc := pc.(models.PC)

	if err = c.auditLog.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ACTION_GET,
		Entity:    audit.ENTITY_COMPUTER,
		EntityID:  int64(gotPc.ID),
		Timestamp: time.Now(),
	}); err != nil {
		config.LogError("getComputer", err)
	}

	return gotPc, nil
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

	if err = c.auditLog.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ACTION_CREATE,
		Entity:    audit.ENTITY_COMPUTER,
		EntityID:  int64(pc.ID),
		Timestamp: time.Now(),
	}); err != nil {
		config.LogError("addComputer", err)
	}

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

	if err = c.auditLog.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ACTION_UPDATE,
		Entity:    audit.ENTITY_COMPUTER,
		EntityID:  int64(pc.ID),
		Timestamp: time.Now(),
	}); err != nil {
		config.LogError("updateComputer", err)
	}

	return nil
}

func (c *ComputersStore) DeleteComputer(ctx context.Context, id int) error {
	err := c.repo.DeleteComputer(ctx, id)
	if err != nil {
		return err
	}

	if err = c.auditLog.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ACTION_DELETE,
		Entity:    audit.ENTITY_COMPUTER,
		EntityID:  int64(id),
		Timestamp: time.Now(),
	}); err != nil {
		config.LogError("deleteComputer", err)
	}

	return nil
}
