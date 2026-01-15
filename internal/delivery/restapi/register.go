package restapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/georgg2003/gophermart/internal/usecase"
	"github.com/labstack/echo/v4"
)

func (s *server) PostAPIUserRegister(c echo.Context) error {
	req := c.Request()
	defer req.Body.Close()
	ctx := req.Context()

	logger := s.logger.WithRequestCtx(ctx)

	decoder := json.NewDecoder(req.Body)

	var registerRequest RegisterRequest
	err := decoder.Decode(&registerRequest)
	if err != nil {
		logger.WithError(err).Error("failed to decode json register request")
		return c.String(http.StatusBadRequest, "wrong request format")
	}

	accessToken, err := s.uc.UserRegister(ctx, registerRequest.Login, registerRequest.Password)
	if errors.Is(err, usecase.ErrUserAlreadyExists) {
		return c.String(http.StatusConflict, "user already exists")
	}
	if err != nil {
		logger.WithError(err).Error("failed to register user")
		return err
	}
	logger.With(slog.String("login", registerRequest.Login)).Info("successfully registered user")

	c.Response().Header().Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", accessToken))
	return c.String(http.StatusOK, "successfully registered")
}
