package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/ninja-way/pc-store/internal/model"
	"github.com/ninja-way/pc-store/internal/repository"
	"time"
)

type PG struct {
	ctx  context.Context
	conn *pgx.Conn
}

func Init(ctx context.Context) (repository.DB, error) {
	psqlInfo := "host=127.0.0.1 port=5432 user=postgres " +
		"password=1234 dbname=pcstore sslmode=disable"

	conn, err := pgx.Connect(ctx, psqlInfo)
	if err != nil {
		return nil, err
	}
	//defer conn.Close(context.Background())

	return &PG{
		ctx:  ctx,
		conn: conn,
	}, nil
}

func (p *PG) GetComputers() ([]model.PC, error) {
	rows, err := p.conn.Query(p.ctx, "select * from pc")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	computers := make([]model.PC, 0)
	for rows.Next() {
		pc := model.PC{}
		// fix nil instead values
		if err = rows.Scan(&pc.ID, &pc.Name, &pc.CPU, nil,
			&pc.RAM, nil, &pc.AddedAt, &pc.Price); err != nil {
			return nil, err
		}
		computers = append(computers, pc)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return computers, nil
}

func (p *PG) GetComputerByID(id int) (model.PC, error) {
	var pc model.PC
	err := p.conn.QueryRow(p.ctx, "select * from pc where id = $1", id).
		Scan(&pc.ID, &pc.Name, &pc.CPU, nil, &pc.RAM, nil, &pc.AddedAt, &pc.Price)
	if err != nil {
		return model.PC{}, err
	}

	return pc, nil
}

func (p *PG) AddComputer(pc model.PC) error {
	if pc.AddedAt.IsZero() {
		pc.AddedAt = time.Now()
	}

	_, err := p.conn.Exec(p.ctx, "insert into pc (name, cpu, videocard, ram, data_storage, added_at, price) values "+
		"($1, $2, $3, $4, $5, $6, $7)", pc.Name, pc.CPU, pc.Videocard, pc.RAM, pc.DataStorage, pc.AddedAt, pc.Price)
	if err != nil {
		return err
	}
	return nil
}

func (p *PG) UpdateComputer(int, model.PC) error {
	return nil
}

func (p *PG) DeleteComputer(int) error {
	return nil
}
