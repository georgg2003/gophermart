package postgres

import (
	"context"

	"github.com/georgg2003/gophermart/internal/pkg/config"
	"github.com/georgg2003/gophermart/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type postgres struct {
	cfg    *config.Config
	logger *logrus.Logger
	db     *pgxpool.Pool
}

func New(cfg *config.Config, logger *logrus.Logger, ctx context.Context) repository.Repository {
	var db *pgxpool.Pool

	logger.WithField("dsn", cfg.DataBaseURI).Debug("making new db connection pool")

	poolConfig, err := pgxpool.ParseConfig(cfg.DataBaseURI)
	if err != nil {
		logger.WithError(err).Fatal("unable to parse DATABASE_URL")
	}

	db, err = pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		logger.WithError(err).Fatal("unable to create connection pool")
	}

	m, err := migrate.New("file://migrations", cfg.DataBaseURI)
	if err != nil {
		logger.WithError(err).Fatal("failed to create migrate instance")
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		logger.WithError(err).Fatal("failed to migrate")
	}
	logger.Info("db successfully migrated")

	return &postgres{
		cfg:    cfg,
		logger: logger,
		db:     db,
	}
}
