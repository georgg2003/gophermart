package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/georgg2003/gophermart/internal/delivery/restapi"
	"github.com/georgg2003/gophermart/internal/pkg/config"
	"github.com/georgg2003/gophermart/internal/pkg/logging"
	"github.com/georgg2003/gophermart/internal/pkg/middleware"
	"github.com/georgg2003/gophermart/internal/repository/accrual"
	"github.com/georgg2003/gophermart/internal/repository/postgres"
	"github.com/georgg2003/gophermart/internal/usecase"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	oapiValidator "github.com/oapi-codegen/echo-middleware"
	"golang.org/x/sync/errgroup"
)

const pathToSwagger = "api/swagger.yaml"

func main() {
	logger := logging.New(slog.LevelDebug)

	cfg := config.New()
	if err := cfg.ReadFromEnv(); err != nil {
		logger.WithError(err).Fatal("failed to read config from env")
	}
	cfg.ReadFromFlags()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	repository, err := postgres.New(cfg, logger.WithString("layer", "pg repo"), ctx)
	if err != nil {
		logger.WithError(err).Fatal("failed to create postgres repository")
	}
	accrualRepo := accrual.New(cfg, logger.WithString("layer", "accrual repo"))

	usecase := usecase.New(
		cfg,
		logger.WithString("layer", "usecase"),
		repository,
		accrualRepo,
	)

	delivery := restapi.NewServer(
		cfg,
		logger.WithString("layer", "delivery"),
		usecase,
	)

	data, err := os.ReadFile(pathToSwagger)
	if err != nil {
		logger.WithError(err).WithString("path", pathToSwagger).Fatal("error reading swagger file")
	}

	swagger, err := openapi3.NewLoader().LoadFromData(data)
	if err != nil {
		logger.WithError(err).WithString("path", pathToSwagger).Fatal("error parsing file as Swagger YAML")
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
		middleware.NewAuthMiddleware(cfg.JWTSecretKey, logger, func(c echo.Context) bool {
			return c.Path() == "/api/user/login" || c.Path() == "/api/user/register"
		}),
		validator,
	)

	restapi.RegisterHandlers(e, delivery)

	g, ctx := errgroup.WithContext(ctx)
	for i := 0; i < cfg.Workers; i++ {
		g.Go(func() error {
			usecase.MakeProcessorWorker(ctx)
			return nil
		})
	}

	g.Go(func() error {
		if err = e.Start(cfg.RunAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.WithError(err).Error("server failed")
			return err
		}
		return nil
	})

	g.Go(func() error {
		<-ctx.Done()
		return e.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil {
		logger.WithError(err).Error("application stopped with error")
	}
}
