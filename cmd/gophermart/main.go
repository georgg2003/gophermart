package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/georgg2003/gophermart/internal/delivery/restapi"
	"github.com/georgg2003/gophermart/internal/pkg/config"
	"github.com/georgg2003/gophermart/internal/repository/postgres"
	"github.com/georgg2003/gophermart/internal/usecase"
	"github.com/labstack/echo/v4"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	repository := postgres.New(cfg, logger)
	usecase := usecase.New(cfg, logger, repository)
	delivery := restapi.NewServer(cfg, logger, usecase)

	restapi.RegisterHandlers(e, delivery)

	// And we serve HTTP until the world ends.
	log.Fatal(e.Start(cfg.RunAddr))
}
