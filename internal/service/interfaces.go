package service

import (
	"context"
	audit "github.com/ninja-way/mq-audit-log/pkg/models"
	"github.com/ninja-way/pc-store/internal/models"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type UsersRepository interface {
	CreateUser(context.Context, models.User) error
	GetUser(context.Context, string, string) (models.User, error)
}

type SessionsRepository interface {
	CreateToken(ctx context.Context, token models.RefreshSession) error
	GetToken(ctx context.Context, token string) (models.RefreshSession, error)
}

type AuditClient interface {
	SendLogRequest(ctx context.Context, req audit.LogItem) error
}

type ComputerRepository interface {
	GetComputers(context.Context) ([]models.PC, error)
	GetComputerByID(context.Context, int) (models.PC, error)
	AddComputer(context.Context, models.PC) (int, error)
	UpdateComputer(context.Context, int, models.PC) error
	DeleteComputer(context.Context, int) error
}
