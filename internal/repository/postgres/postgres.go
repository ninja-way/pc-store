package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/ninja-way/pc-store/internal/model"
	"github.com/ninja-way/pc-store/internal/repository"
)

type PG struct {
	conn *pgx.Conn
}

func Init() (repository.DB, error) {
	psqlInfo := "host=127.0.0.1 port=5432 user=postgres " +
		"password=1234 dbname=pcstore sslmode=disable"

	conn, err := pgx.Connect(context.Background(), psqlInfo)
	if err != nil {
		return nil, err
	}
	//defer conn.Close(context.Background())

	return PG{conn: conn}, nil
}

func (p PG) GetComputers() ([]model.PC, error) {
	rows, err := p.conn.Query(context.Background(), "select * from pc")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	computers := make([]model.PC, 0)
	for rows.Next() {
		pc := model.PC{}
		// fix it
		if err = rows.Scan(&pc.ID, &pc.Name, &pc.CPU, nil, &pc.RAM, nil, nil, &pc.Price); err != nil {
			return nil, err
		}
		computers = append(computers, pc)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	fmt.Println(computers)

	return computers, nil
}

func (p PG) GetComputerByID(int) (model.PC, error) {
	return model.PC{}, nil
}

func (p PG) AddComputer(model.PC) error {
	return nil
}

func (p PG) UpdateComputer(int, model.PC) error {
	return nil
}

func (p PG) DeleteComputer(int) error {
	return nil
}
