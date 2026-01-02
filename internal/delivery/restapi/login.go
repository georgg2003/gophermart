package restapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/georgg2003/gophermart/internal/usecase"
	"github.com/labstack/echo/v4"
)

func (s *server) PostApiUserLogin(c echo.Context) error {
	req := c.Request()
	defer req.Body.Close()
	ctx := req.Context()

	decoder := json.NewDecoder(req.Body)

	var loginRequest LoginRequest
	err := decoder.Decode(&loginRequest)
	if err != nil {
		s.logger.WithError(err).Error("failed to decode json login request")
		return c.String(http.StatusBadRequest, "wrong request format")
	}

	accessToken, err := s.uc.UserLogin(ctx, loginRequest.Login, loginRequest.Password)
	if err != nil {
		if errors.Is(err, usecase.ErrUserWrongPassword) ||
			errors.Is(err, usecase.ErrUserNotFound) {
			return c.String(http.StatusUnauthorized, "Incorrect login or password")
		}
		s.logger.WithError(err).Error("failed to login user")
		return err
	}
	s.logger.WithField("login", loginRequest.Login).Info("successfully logged in user")

	c.Response().Header().Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", accessToken))
	c.Response().WriteHeader(http.StatusOK)

	return nil
}
