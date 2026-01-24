package domain

import (
	"context"
	"sellers-accounts-backend/internal/entities"
)

func (u *Usecase) Deals(ctx context.Context, req *entities.RequestGetDeals) (*entities.ResponseGetDeals, error) {
	deals, total, err := u.repo.Deals(ctx, req.Limit, req.Page)
	if err != nil {
		u.log.Errorw("usecase.AllAccounts", "error", err)
		return nil, mapRepoError(err)
	}
	return &entities.ResponseGetDeals{
		Deals: deals,
		Total: total,
	}, nil
}
