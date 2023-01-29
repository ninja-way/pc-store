package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/ninja-way/pc-store/internal/model"
	"github.com/ninja-way/pc-store/internal/repository"
	"time"
)

// db config
const (
	host     = "127.0.0.1"
	port     = "5432"
	user     = "postgres"
	password = 1234
	dbname   = "pcstore"
)

// PG implements DB interface
type PG struct {
	ctx  context.Context
	conn *pgx.Conn
}

// Init makes connection to the database with the passed context
func Init(ctx context.Context) (repository.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%d dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	conn, err := pgx.Connect(ctx, psqlInfo)
	if err != nil {
		return nil, err
	}

	return &PG{
		ctx:  ctx,
		conn: conn,
	}, nil
}

// Close postgres db connection
func (p *PG) Close(ctx context.Context) error {
	return p.conn.Close(ctx)
}

// GetComputers return all pc from db
func (p *PG) GetComputers() ([]model.PC, error) {
	rows, err := p.conn.Query(p.ctx, "select * from pc")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	computers := make([]model.PC, 0)
	for rows.Next() {
		pc := model.PC{}
		err = rows.Scan(&pc.ID, &pc.Name, &pc.CPU, &pc.Videocard, &pc.RAM,
			&pc.DataStorage, &pc.AddedAt, &pc.Price)
		if err != nil {
			return nil, err
		}
		computers = append(computers, pc)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return computers, nil
}

// GetComputerByID return pc by id from db
func (p *PG) GetComputerByID(id int) (model.PC, error) {
	var pc model.PC

	err := p.conn.QueryRow(p.ctx, "select * from pc where id = $1", id).
		Scan(&pc.ID, &pc.Name, &pc.CPU, &pc.Videocard, &pc.RAM, &pc.DataStorage, &pc.AddedAt, &pc.Price)
	if err != nil {
		return model.PC{}, err
	}

	return pc, nil
}

// AddComputer insert passed pc into db
// if AddedAt value not specified, sets the current time
func (p *PG) AddComputer(pc model.PC) error {
	if pc.AddedAt.IsZero() {
		pc.AddedAt = time.Now()
	}

	_, err := p.conn.Exec(p.ctx, "insert into pc (name, cpu, videocard, ram, data_storage, added_at, price) "+
		"values ($1, $2, $3, $4, $5, $6, $7)",
		pc.Name, pc.CPU, pc.Videocard, pc.RAM, pc.DataStorage, pc.AddedAt, pc.Price)

	return err
}

// UpdateComputer changes only the specified fields in the PC in computer with passed id
func (p *PG) UpdateComputer(id int, newPC model.PC) error {
	var newParam = make(map[string]interface{})

	if newPC.Name != "" {
		newParam["name"] = newPC.Name
	}
	if newPC.CPU != "" {
		newParam["cpu"] = newPC.CPU
	}
	if newPC.Videocard != "" {
		newParam["videocard"] = newPC.Videocard
	}
	if newPC.RAM != 0 {
		newParam["ram"] = newPC.RAM
	}
	if newPC.DataStorage != "" {
		newParam["data_storage"] = newPC.DataStorage
	}
	if newPC.Price != 0 {
		newParam["price"] = newPC.Price
	}

	for i, v := range newParam {
		if _, err := p.conn.Exec(p.ctx, "update pc set "+i+"=$1 where id = $2", v, id); err != nil {
			return err
		}
	}
	return nil
}

// DeleteComputer from db by id
func (p *PG) DeleteComputer(id int) error {
	t, err := p.conn.Exec(p.ctx, "delete from pc where id = $1", id)
	if t.RowsAffected() == 0 {
		return errors.New("no rows in result set")
	}
	return err
}
