package restapi

import (
	"errors"
	"net/http"

	"github.com/georgg2003/gophermart/internal/pkg/contextlib"
	"github.com/georgg2003/gophermart/internal/usecase"
	"github.com/labstack/echo/v4"
)

func (s *server) GetApiUserWithdrawals(c echo.Context) error {
	req := c.Request()
	defer req.Body.Close()
	ctx := req.Context()

	userID := contextlib.MustGetUserID(ctx, s.logger)
	withdrawls, err := s.uc.UserGetWithdrawals(ctx, userID)
	if err != nil {
		if errors.Is(err, usecase.ErrWidthdrawalsNotFound) {
			return c.String(http.StatusNoContent, "No withdrawals found")
		}
		s.logger.WithError(err).Error("failed to get user balance")
		return err
	}

	withdrawalsDTO := make([]WithdrawalInfo, 0, len(withdrawls))
	for _, v := range withdrawls {
		withdrawalsDTO = append(withdrawalsDTO, WithdrawalInfo{
			Order:       v.Order,
			Sum:         v.Sum.Major(),
			ProcessedAt: v.ProcessedAt,
		})
	}
	c.JSON(http.StatusOK, withdrawalsDTO)

	return nil
}
