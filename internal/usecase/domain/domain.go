package domain

import (
	"context"
	"sellers-accounts-backend/internal/repository"

	"go.uber.org/zap"
	"sellers-accounts-backend/config"
)

type Usecase struct {
	cfg  *config.Config
	log  *zap.SugaredLogger
	ctx  context.Context
	repo repository.InterfaceRepository
}

func New(cfg *config.Config, log *zap.SugaredLogger, ctx context.Context, repo repository.InterfaceRepository) *Usecase {
	return &Usecase{
		cfg:  cfg,
		log:  log,
		ctx:  ctx,
		repo: repo,
	}
}
