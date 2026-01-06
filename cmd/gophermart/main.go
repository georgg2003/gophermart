package main

import (
	"context"
	"os"

	"github.com/georgg2003/gophermart/internal/delivery/restapi"
	"github.com/georgg2003/gophermart/internal/pkg/config"
	"github.com/georgg2003/gophermart/internal/pkg/middleware"
	"github.com/georgg2003/gophermart/internal/repository/accrual"
	"github.com/georgg2003/gophermart/internal/repository/postgres"
	"github.com/georgg2003/gophermart/internal/usecase"
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
		logger.WithError(err).Fatal("failed to create config")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repository, err := postgres.New(cfg, logger, ctx)
	if err != nil {
		logger.WithError(err).Fatal("failed to create postgres repository")
	}
	accrualRepo := accrual.New(cfg)

	usecase := usecase.New(cfg, logger, repository, accrualRepo)
	delivery := restapi.NewServer(cfg, logger, usecase)

	data, err := os.ReadFile(pathToSwagger)
	if err != nil {
		logger.WithError(err).Fatalf("error reading %s", pathToSwagger)
	}

	swagger, err := openapi3.NewLoader().LoadFromData(data)
	if err != nil {
		logger.WithError(err).Fatalf("error parsing %s as Swagger YAML", pathToSwagger)
	}

	validator := oapiValidator.OapiRequestValidatorWithOptions(
		swagger,
		&oapiValidator.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: openapi3filter.NoopAuthenticationFunc,
			},
		},
	)

	if err != nil {
		logger.WithError(err).Fatal("failed to create oapi validator")
	}

	e := echo.New()
	e.Use(
		echoMiddleware.RequestID(),
		echoMiddleware.Decompress(),
		echoMiddleware.Gzip(),
		echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		}),
		echoMiddleware.Recover(),
		middleware.LoggingMiddleware(logger),
		middleware.NewAuthMiddleware(cfg, logger, func(c echo.Context) bool {
			if c.Path() == "/api/user/login" || c.Path() == "/api/user/register" {
				return true
			}
			return false
		}),
		validator,
	)

	restapi.RegisterHandlers(e, delivery)

	for i := 0; i < cfg.Workers; i++ {
		go usecase.MakeProcessorWorker(ctx)
	}

	err = e.Start(cfg.RunAddr)
	if err != nil {
		logger.WithError(err).Fatal("server failed")
	}
}
