package postgres

import (
	"context"
	"github.com/ninja-way/pc-store/internal/models"
)

func (p *PG) CreateUser(ctx context.Context, user models.User) error {
	_, err := p.conn.Exec(ctx, "insert into users (name, email, password, registered_at) values ($1, $2, $3, $4)",
		user.Name, user.Email, user.Password, user.RegisteredAt)
	return err
}
