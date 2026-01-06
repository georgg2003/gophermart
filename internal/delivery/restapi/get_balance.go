package restapi

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *server) GetAPIUserBalance(c echo.Context) error {
	req := c.Request()
	defer req.Body.Close()
	ctx := req.Context()

	balance, err := s.uc.UserGetBalance(ctx)
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
