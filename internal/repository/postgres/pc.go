package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/ninja-way/pc-store/internal/config"
	"github.com/ninja-way/pc-store/internal/models"
	"github.com/ninja-way/pc-store/internal/repository"
)

// PG is postgres connection implements DB interface
type PG struct {
	conn *pgx.Conn
}

// Connect makes connection to the database with the passed context
func Connect(ctx context.Context, db *config.Postgres) (repository.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		db.Host, db.Port, db.UserName, db.Password, db.DBName, db.SSLMode)

	conn, err := pgx.Connect(ctx, psqlInfo)
	if err != nil {
		return nil, err
	}

	return &PG{conn: conn}, nil
}

// Close postgres db connection
func (p *PG) Close(ctx context.Context) error {
	return p.conn.Close(ctx)
}

// GetComputers return all pc from db
func (p *PG) GetComputers(ctx context.Context) ([]models.PC, error) {
	rows, err := p.conn.Query(ctx, "select * from pc")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	computers := make([]models.PC, 0)
	for rows.Next() {
		pc := models.PC{}
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
func (p *PG) GetComputerByID(ctx context.Context, id int) (models.PC, error) {
	var pc models.PC

	err := p.conn.QueryRow(ctx, "select * from pc where id = $1", id).
		Scan(&pc.ID, &pc.Name, &pc.CPU, &pc.Videocard, &pc.RAM, &pc.DataStorage, &pc.AddedAt, &pc.Price)
	if err != nil {
		return models.PC{}, err
	}

	return pc, nil
}

// AddComputer insert passed pc into db
func (p *PG) AddComputer(ctx context.Context, pc models.PC) (int, error) {
	var id int

	var insertQuery = "insert into pc (name, cpu, videocard, ram, data_storage, added_at, price) values ($1, $2, $3, $4, $5, $6, $7) returning id"
	row := p.conn.QueryRow(ctx, insertQuery, pc.Name, pc.CPU, pc.Videocard, pc.RAM, pc.DataStorage, pc.AddedAt, pc.Price)

	err := row.Scan(&id)
	return id, err
}

// UpdateComputer changes only the specified fields in the PC in computer with passed id
func (p *PG) UpdateComputer(ctx context.Context, id int, newPC models.PC) error {
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

	for param, val := range newParam {
		if _, err := p.conn.Exec(ctx, "update pc set "+param+"=$1 where id = $2", val, id); err != nil {
			return err
		}
	}
	return nil
}

// DeleteComputer from db by id
func (p *PG) DeleteComputer(ctx context.Context, id int) error {
	t, err := p.conn.Exec(ctx, "delete from pc where id = $1", id)
	if t.RowsAffected() == 0 {
		return errors.New("id not found")
	}
	return err
}
