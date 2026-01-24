package domain

import (
	"context"
	"sellers-accounts-backend/internal/entities"
)

func (u *Usecase) AccountData(ctx context.Context, req *entities.RequestGetAccountData) (*entities.ResponseGetAccountData, error) {
	if req.ID == 0 {
		u.log.Errorw("account data id invalid", "req", req)
		return nil, ErrInvalidParams
	}
	data, err := u.repo.AccountData(ctx, req.ID)
	if err != nil {
		u.log.Errorw("usecase.AccountData", "error", err, "id", req.ID)
		return nil, mapRepoError(err)
	}
	return &entities.ResponseGetAccountData{Data: *data}, nil
}

func (u *Usecase) DeleteAccountData(ctx context.Context, req *entities.RequestDeleteAccountData) (*entities.ResponseDeleteAccountData, error) {
	if req.ID == 0 {
		u.log.Errorw("account data id invalid", "req", req)
		return nil, ErrInvalidParams
	}
	err := u.repo.DeleteAccountData(ctx, req.ID)
	if err != nil {
		u.log.Errorw("usecase.DeleteAccountData", "error", err, "id", req.ID)
		return nil, mapRepoError(err)
	}
	return &entities.ResponseDeleteAccountData{Status: true}, nil
}
