package restapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/georgg2003/gophermart/internal/usecase"
	"github.com/labstack/echo/v4"
)

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

	c.Response().Header().Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", accessToken))
	c.Response().WriteHeader(http.StatusOK)
	return nil
}
