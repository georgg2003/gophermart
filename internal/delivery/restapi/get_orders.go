package restapi

import (
	"errors"
	"net/http"

	"github.com/georgg2003/gophermart/internal/usecase"
	"github.com/labstack/echo/v4"
)

func (s *server) GetAPIUserOrders(c echo.Context) error {
	req := c.Request()
	defer req.Body.Close()
	ctx := req.Context()

	orders, err := s.uc.UserGetOrders(ctx)
	if err != nil {
		if errors.Is(err, usecase.ErrOrdersNotFound) {
			return c.String(http.StatusNoContent, "No orders found")
		}
		s.logger.WithError(err).Error("failed to get user orders")
		return err
	}

	ordersDTO := make([]OrderInfo, 0, len(orders))
	for _, v := range orders {
		ordersDTO = append(ordersDTO, OrderInfo{
			Accrual:    v.Accrual.NullMajor(),
			Number:     v.Number,
			Status:     OrderInfoStatus(v.Status),
			UploadedAt: v.UploadedAt,
		})
	}
	return c.JSON(http.StatusOK, ordersDTO)
}
