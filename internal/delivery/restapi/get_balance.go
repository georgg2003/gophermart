package restapi

import (
	"net/http"

	"github.com/georgg2003/gophermart/internal/pkg/contextlib"
	"github.com/labstack/echo/v4"
)

func (s *server) GetApiUserBalance(c echo.Context) error {
	req := c.Request()
	defer req.Body.Close()
	ctx := req.Context()

	userID := contextlib.MustGetUserID(ctx, s.logger)
	balance, err := s.uc.UserGetBalance(ctx, userID)
	if err != nil {
		s.logger.WithError(err).Error("failed to get user balance")
		return err
	}
	resp := BalanceResponse{
		Current:   balance.Current.Major(),
		Withdrawn: balance.Withdrawn.Major(),
	}
	c.JSON(http.StatusOK, resp)

	return nil
}
