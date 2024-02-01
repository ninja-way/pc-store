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

func (p *PG) GetUser(ctx context.Context, email, password string) (models.User, error) {
	var user models.User
	err := p.conn.QueryRow(ctx, "select id, name, email, registered_at from users where email = $1 and password = $2",
		email, password).Scan(&user.ID, &user.Name, &user.Email, &user.RegisteredAt)
	return user, err
}
