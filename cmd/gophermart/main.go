package main

import (
	"context"
	"os"

	"github.com/georgg2003/gophermart/internal/delivery/restapi"
	"github.com/georgg2003/gophermart/internal/pkg/config"
	"github.com/georgg2003/gophermart/internal/repository/postgres"
	"github.com/georgg2003/gophermart/internal/usecase"
	"github.com/georgg2003/gophermart/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	oapiValidator "github.com/oapi-codegen/echo-middleware"
	"github.com/sirupsen/logrus"
)

const pathToSwagger = "api/swagger.yaml"

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

	repository, err := postgres.New(cfg, logger, ctx)
	if err != nil {
		logger.WithError(err).Fatal("failed to create postgres repository")
	}

	usecase := usecase.New(cfg, logger, repository)
	delivery := restapi.NewServer(cfg, logger, usecase)

	data, err := os.ReadFile(pathToSwagger)
	if err != nil {
		logger.WithError(err).Fatal("error reading %s", pathToSwagger)
	}

	swagger, err := openapi3.NewLoader().LoadFromData(data)
	if err != nil {
		logger.WithError(err).Fatal("error parsing %s as Swagger YAML", pathToSwagger)
	}

	validator := oapiValidator.OapiRequestValidatorWithOptions(
		swagger,
		&oapiValidator.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: func(context.Context, *openapi3filter.AuthenticationInput) error {
					return nil // TODO
				},
			},
		},
	)

	if err != nil {
		logger.WithError(err).Fatal("failed to create oapi validator")
	}

	e := echo.New()
	e.Use(
		echoMiddleware.RequestID(),
		middleware.LoggingMiddleware(logger),
		echoMiddleware.Gzip(),
		echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		}),
		validator,
	)

	restapi.RegisterHandlers(e, delivery)

	err = e.Start(cfg.RunAddr)
	if err != nil {
		logger.WithError(err).Fatal("server failed")
	}
}
