package usecase

import (
	"github.com/georgg2003/gophermart/internal/pkg/config"
	"github.com/georgg2003/gophermart/internal/repository"
	"github.com/sirupsen/logrus"
)

type useCase struct {
	cfg    *config.Config
	logger *logrus.Logger
	repo   repository.Repository
}

type UseCase interface {
}

func New(
	cfg *config.Config,
	logger *logrus.Logger,
	repo repository.Repository,
) UseCase {
	return &useCase{
		cfg:    cfg,
		logger: logger,
		repo:   repo,
	}
}
