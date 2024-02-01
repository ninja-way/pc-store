package postgres

import (
	"context"
	"github.com/ninja-way/pc-store/internal/models"
)

func (p *PG) CreateToken(ctx context.Context, token models.RefreshSession) error {
	_, err := p.conn.Exec(ctx, "INSERT INTO refresh_tokens (user_id, token, expires_at) values ($1, $2, $3)",
		token.UserID, token.Token, token.ExpiresAt)

	return err
}

func (p *PG) GetToken(ctx context.Context, token string) (models.RefreshSession, error) {
	var t models.RefreshSession
	err := p.conn.QueryRow(ctx, "SELECT id, user_id, token, expires_at FROM refresh_tokens WHERE token=$1", token).
		Scan(&t.ID, &t.UserID, &t.Token, &t.ExpiresAt)
	if err != nil {
		return t, err
	}

	_, err = p.conn.Exec(ctx, "DELETE FROM refresh_tokens WHERE user_id = $1", t.UserID)

	return t, err
}
