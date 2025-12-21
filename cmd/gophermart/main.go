package main

import (
	"context"

	"github.com/georgg2003/gophermart/internal/delivery/restapi"
	"github.com/georgg2003/gophermart/internal/pkg/config"
	"github.com/georgg2003/gophermart/internal/repository/postgres"
	"github.com/georgg2003/gophermart/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.DebugLevel)

	cfg, err := config.New()
	if err != nil {
		logger.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repository := postgres.New(cfg, logger, ctx)
	usecase := usecase.New(cfg, logger, repository)
	delivery := restapi.NewServer(cfg, logger, usecase)

	e := echo.New()
	e.Use(middleware.Gzip())

	restapi.RegisterHandlers(e, delivery)

	logger.Fatal(e.Start(cfg.RunAddr))
}
