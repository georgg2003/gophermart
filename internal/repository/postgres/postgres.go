package postgres

import (
	"context"
	"log/slog"

	"github.com/georgg2003/gophermart/internal/pkg/config"
	"github.com/georgg2003/gophermart/internal/pkg/logging"
	"github.com/georgg2003/gophermart/internal/repository"
	"github.com/georgg2003/gophermart/pkg/errutils"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type postgres struct {
	cfg    *config.Config
	logger *logging.Logger
	db     *pgxpool.Pool
}

func New(cfg *config.Config, logger *logging.Logger, ctx context.Context) (repository.Repository, error) {
	var db *pgxpool.Pool

	logger.With(slog.String("dsn", cfg.DataBaseURI)).Debug("making new db connection pool")

	poolConfig, err := pgxpool.ParseConfig(cfg.DataBaseURI)
	if err != nil {
		return nil, errutils.Wrap(err, "unable to parse DATABASE_URL")
	}

	db, err = pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, errutils.Wrap(err, "unable to create connection pool")
	}

	m, err := migrate.New("file://migrations", cfg.DataBaseURI)
	if err != nil {
		return nil, errutils.Wrap(err, "failed to create migrate instance")
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return nil, errutils.Wrap(err, "failed apply migration")
	}
	logger.Info("db successfully migrated")

	return &postgres{
		cfg:    cfg,
		logger: logger,
		db:     db,
	}, nil
}
