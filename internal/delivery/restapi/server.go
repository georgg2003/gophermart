package restapi

import (
	"github.com/georgg2003/gophermart/internal/pkg/config"
	"github.com/georgg2003/gophermart/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

//go:generate go tool oapi-codegen --config=./gen/server.yaml ../../../api/swagger.yaml
//go:generate go tool oapi-codegen --config=./gen/models.yaml ../../../api/swagger.yaml

type server struct {
	cfg    *config.Config
	logger *logrus.Logger
	uc     usecase.UseCase
}

func NewServer(
	cfg *config.Config,
	logger *logrus.Logger,
	uc usecase.UseCase,
) ServerInterface {
	return &server{
		cfg:    cfg,
		logger: logger,
		uc:     uc,
	}
}

func (s *server) PostApiUserBalanceWithdraw(ctx echo.Context) error { return nil }
func (s *server) GetApiUserOrders(ctx echo.Context) error           { return nil }
func (s *server) PostApiUserOrders(ctx echo.Context) error          { return nil }
func (s *server) GetApiUserWithdrawals(ctx echo.Context) error      { return nil }
