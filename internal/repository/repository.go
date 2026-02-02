package repository

import (
	"context"
	"fmt"
	"sellers-accounts-backend/internal/entities"

	"sellers-accounts-backend/config"

	"sellers-accounts-backend/internal/repository/postgres"

	"go.uber.org/zap"
)

type InterfaceLifecycle interface {
	OnStart(_ context.Context) error
	OnStop(_ context.Context) error
}

type ImplementationRepository interface {
	AllAccounts(ctx context.Context, limit int, page int) (entities.Accounts, int64, error)
	AcceptableTypesAccounts(ctx context.Context) (entities.AcceptableTypesAccounts, error)
	Account(ctx context.Context, uid string) (*entities.Account, error)
	CreateAccount(ctx context.Context, a *entities.Account) (*entities.Account, error)
	UpdateAccount(ctx context.Context, a *entities.Account) error
	DeactivateAccountsByName(ctx context.Context, name string) (int, error)
	DeleteAccount(ctx context.Context, uid string) error
	AccountData(ctx context.Context, id int) (*entities.Data, error)
	DeleteAccountData(ctx context.Context, id int) error
	Deals(ctx context.Context, limit int, page int) (*[]entities.Deal, int64, error)
	BotLink(ctx context.Context) (string, error)
}

type InterfaceRepository interface {
	InterfaceLifecycle
	ImplementationRepository
}

func New(name string, log *zap.SugaredLogger, cfg *config.Config, ctx context.Context) (InterfaceRepository, error) {
	switch name {
	case "postgres":
		return postgres.New(log, cfg, ctx), nil
	default:
		return nil, fmt.Errorf("unknown repo backend: %s", name)
	}
}
