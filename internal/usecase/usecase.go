package usecase

import (
	"context"
	"go.uber.org/zap"
	"sellers-accounts-backend/config"
	"sellers-accounts-backend/internal/entities"
	"sellers-accounts-backend/internal/repository"
	"sellers-accounts-backend/internal/usecase/domain"
)

type InterfaceUsecase interface {
	AllAccounts(ctx context.Context, req *entities.RequestAllAccounts) (*entities.ResponseAllAccounts, error)
	AcceptableTypesAccounts(ctx context.Context) (*entities.ResponseAcceptableTypeAccounts, error)
	Account(ctx context.Context, req *entities.RequestAccountByID) (*entities.ResponseAccountByID, error)
	// CreateAccounts CreateAccount(ctx context.Context, req *entities.RequestCreateAccounts) (*entities.ResponseCreateAccounts, error)
	CreateAccounts(ctx context.Context, req *entities.RequestCreateAccounts) (*entities.ResponseCreateAccounts, error)
	UpdateAccount(ctx context.Context, req *entities.RequestUpdateAccount) (*entities.ResponseUpdateAccount, error)
	DeactivateAccountsByName(ctx context.Context, req *entities.RequestDeactivateAccountsByName) (*entities.ResponseDeactivateAccounts, error)
	AccountData(ctx context.Context, req *entities.RequestGetAccountData) (*entities.ResponseGetAccountData, error)
	DeleteAccountData(ctx context.Context, req *entities.RequestDeleteAccountData) (*entities.ResponseDeleteAccountData, error)
	Deals(ctx context.Context, req *entities.RequestGetDeals) (*entities.ResponseGetDeals, error)
}

func New(cfg *config.Config, log *zap.SugaredLogger, cxt context.Context, repo repository.InterfaceRepository) InterfaceUsecase {
	return domain.New(cfg, log, cxt, repo)
}
