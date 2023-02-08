package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/ninja-way/pc-store/internal/config"
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
