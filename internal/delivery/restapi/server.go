package restapi

import (
	"encoding/json"
	"errors"
	"net/http"

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

func (s *server) PostApiUserRegister(c echo.Context) error {
	req := c.Request()
	defer req.Body.Close()
	ctx := req.Context()

	decoder := json.NewDecoder(req.Body)

	var registerRequest RegisterRequest
	err := decoder.Decode(&registerRequest)
	if err != nil {
		s.logger.WithError(err).Error("failed to decode json register request")
		return c.String(http.StatusBadRequest, "wrong request format")
	}

	accessToken, err := s.uc.UserRegister(ctx, registerRequest.Login, registerRequest.Password)
	if errors.Is(err, usecase.ErrUserAlreadyExists) {
		return c.String(http.StatusConflict, "user already exists")
	}
	if err != nil {
		s.logger.WithError(err).Error("failed to register user")
		return err
	}
	s.logger.WithField("login", registerRequest.Login).Info("successfully registered user")

	c.Response().Header().Set(echo.HeaderAuthorization, accessToken)
	c.Response().WriteHeader(http.StatusOK)
	return nil
}

func (s *server) PostApiUserLogin(ctx echo.Context) error { return nil }

func (s *server) GetApiUserBalance(ctx echo.Context) error          { return nil }
func (s *server) PostApiUserBalanceWithdraw(ctx echo.Context) error { return nil }
func (s *server) GetApiUserOrders(ctx echo.Context) error           { return nil }
func (s *server) PostApiUserOrders(ctx echo.Context) error          { return nil }
func (s *server) GetApiUserWithdrawals(ctx echo.Context) error      { return nil }
