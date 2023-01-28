package cache

import (
	"errors"
	"github.com/ninja-way/pc-store/internal/model"
	"github.com/ninja-way/pc-store/internal/repository"
	"time"
)

type Cache struct {
	data *[]model.PC
}

func Init() repository.DB {
	return Cache{data: &[]model.PC{}}
}

func (c Cache) GetComputers() ([]model.PC, error) {
	return *c.data, nil
}

func (c Cache) GetComputerByID(id int) (model.PC, error) {
	for _, pc := range *c.data {
		if pc.ID == id {
			return pc, nil
		}
	}
	return model.PC{}, errors.New("pc not found")
}

func (c Cache) AddComputer(pc model.PC) error {
	pc.ID = len(*c.data) + 1
	pc.AddedAt = time.Now()
	*c.data = append(*c.data, pc)
	return nil
}

func (c Cache) UpdateComputer(id int, newPc model.PC) error {
	for i, pc := range *c.data {
		if pc.ID == id {
			newPc.ID = id

			temp := *c.data
			temp[i] = newPc
			c.data = &temp
			return nil
		}
	}
	return errors.New("pc not found")
}

func (c Cache) DeleteComputer(id int) error {
	for i, pc := range *c.data {
		if pc.ID == id {
			temp := *c.data
			*c.data = append(temp[:i], temp[i+1:]...)
			return nil
		}
	}
	return errors.New("pc not found")
}
