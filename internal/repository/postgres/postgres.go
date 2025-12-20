package postgres

import (
	"log/slog"

	"github.com/georgg2003/gophermart/internal/pkg/config"
	"github.com/georgg2003/gophermart/internal/repository"
)

type postgres struct {
	cfg    *config.Config
	logger *slog.Logger
}

func New(cfg *config.Config, logger *slog.Logger) repository.Repository {
	return &postgres{
		cfg:    cfg,
		logger: logger,
	}
}
