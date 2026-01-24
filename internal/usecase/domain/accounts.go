package domain

import (
	"context"
	"errors"
	"sellers-accounts-backend/internal/entities"

	"sellers-accounts-backend/internal/repository/postgres"
)

func (u *Usecase) AllAccounts(ctx context.Context, req *entities.RequestAllAccounts) (*entities.ResponseAllAccounts, error) {
	if req.Limit <= 0 {
		req.Limit = 50
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	accounts, total, err := u.repo.AllAccounts(ctx, req.Limit, req.Page)
	if err != nil {
		u.log.Errorw("usecase.AllAccounts", "error", err)
		return nil, mapRepoError(err)
	}

	//responseAccounts := make([]*entities.Account, 0, len(*accounts))
	//for i := range *accounts {
	//	responseAccounts = append(responseAccounts, &(*accounts)[i])
	//}

	return &entities.ResponseAllAccounts{
		Accounts: accounts,
		Total:    total,
	}, nil
}

func (u *Usecase) AcceptableTypesAccounts(ctx context.Context) (*entities.ResponseAcceptableTypeAccounts, error) {
	types, err := u.repo.AcceptableTypesAccounts(ctx)
	if err != nil {
		u.log.Errorw("usecase.AcceptableTypesAccounts", "error", err)
		return nil, mapRepoError(err)
	}
	if types == nil {
		return nil, nil
	}

	return &entities.ResponseAcceptableTypeAccounts{Types: types}, nil
}

func (u *Usecase) Account(ctx context.Context, req *entities.RequestAccountByID) (*entities.ResponseAccountByID, error) {
	if req.ID == "" {
		return nil, ErrInvalidParams
	}
	account, err := u.repo.Account(ctx, req.ID)
	if err != nil {
		u.log.Errorw("usecase.Account", "error", err, "uid", req.ID)
		return nil, mapRepoError(err)
	}
	return &entities.ResponseAccountByID{Account: account}, nil
}

func (u *Usecase) CreateAccounts(ctx context.Context, req *entities.RequestCreateAccounts) (*entities.ResponseCreateAccounts, error) {
	if req == nil || len(req.Accounts) == 0 {
		return nil, ErrInvalidParams
	}

	created := make([]*entities.Account, 0, len(req.Accounts))
	failed := make([]entities.AccountCreateFailure, 0)

	for i := range req.Accounts {
		item := req.Accounts[i]
		acc, err := u.repo.CreateAccount(ctx, item.ToAccount())
		if err != nil {
			failed = append(failed, entities.AccountCreateFailure{
				Index:   i,
				ID:      item.UID,
				Message: mapRepoError(err).Error(),
			})
			continue
		}
		created = append(created, acc)
	}

	return &entities.ResponseCreateAccounts{
		Accounts: created,
		Failed:   failed,
		Total:    len(req.Accounts),
	}, nil
}

func (u *Usecase) UpdateAccount(ctx context.Context, req *entities.RequestUpdateAccount) (*entities.ResponseUpdateAccount, error) {
	if req.UID == "" {
		return nil, ErrInvalidAccount
	}
	acc := req.ToAccount()
	if err := u.repo.UpdateAccount(ctx, acc); err != nil {
		u.log.Infow("usecase.UpdateAccount", "error", err, "uid", req.UID)
		return nil, mapRepoError(err)
	}
	updated, err := u.repo.Account(ctx, req.UID)
	if err != nil {
		u.log.Infow("usecase.UpdateAccount - fetch account", "error", err, "uid", req.UID)
		return nil, mapRepoError(err)
	}

	return &entities.ResponseUpdateAccount{Account: updated.ToUpdate()}, nil
}

func (u *Usecase) DeactivateAccountsByName(ctx context.Context, req *entities.RequestDeactivateAccountsByName) (*entities.ResponseDeactivateAccounts, error) {
	if req == nil || req.Name == "" {
		return nil, ErrInvalidParams
	}
	updated, err := u.repo.DeactivateAccountsByName(ctx, req.Name)
	if err != nil {
		u.log.Infow("usecase.DeactivateAccountsByName", "error", err, "name", req.Name)
		return nil, mapRepoError(err)
	}

	return &entities.ResponseDeactivateAccounts{Updated: updated}, nil
}

func (u *Usecase) DeleteAccount(ctx context.Context, req *entities.RequestDeleteAccount) error {
	if req.ID == "" {
		return ErrInvalidParams
	}
	if err := u.repo.DeleteAccount(ctx, req.ID); err != nil {
		u.log.Errorw("usecase.DeleteAccount", "error", err, "uid", req.ID)
		return mapRepoError(err)
	}
	return nil
}

func mapRepoError(err error) error {
	switch {
	case errors.Is(err, postgres.ErrNotFoundAccount):
		return ErrNotFoundAccount
	case errors.Is(err, postgres.ErrNilAccount):
		return ErrNilAccount
	case errors.Is(err, postgres.ErrGetAccounts),
		errors.Is(err, postgres.ErrScanAccount):
		return ErrGetAccounts
	case errors.Is(err, postgres.ErrGetAccount):
		return ErrGetAccount
	case errors.Is(err, postgres.ErrGetTypes),
		errors.Is(err, postgres.ErrScanType):
		return ErrGetTypes
	case errors.Is(err, postgres.ErrCreateAccount):
		return ErrCreateAccount
	case errors.Is(err, postgres.ErrUpdateAccount):
		return ErrUpdateAccount
	case errors.Is(err, postgres.ErrDeleteAccount):
		return ErrDeleteAccount
	case errors.Is(err, postgres.ErrDeleteAccountData):
		return ErrDeleteAccountData
	case errors.Is(err, postgres.ErrNotFoundAccountData):
		return ErrNotFoundAccountData
	case errors.Is(err, postgres.ErrGetAccountData):
		return ErrGetAccountData
	case errors.Is(err, postgres.ErrScanDeals), errors.Is(err, postgres.ErrGetDeals):
		return ErrGetDeals
	default:
		return ErrInternal
	}
}
