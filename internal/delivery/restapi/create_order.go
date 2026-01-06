package restapi

import (
	"errors"
	"io"
	"net/http"

	"github.com/georgg2003/gophermart/internal/usecase"
	"github.com/georgg2003/gophermart/pkg/luhn"
	"github.com/labstack/echo/v4"
)

func (s *server) PostAPIUserOrders(c echo.Context) error {
	req := c.Request()
	defer req.Body.Close()
	ctx := req.Context()

	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		s.logger.WithContext(ctx).WithError(err).Error("failed to read body")
		return c.String(http.StatusBadRequest, "wrong body format")
	}
	orderNumber := string(bodyBytes)

	if !luhn.ValidLuhn(orderNumber) {
		s.logger.WithContext(ctx).WithError(err).Error("order number is not valid")
		return c.String(http.StatusUnprocessableEntity, "order number is not valid")
	}

	err = s.uc.UserCreateOrder(ctx, orderNumber)
	if err != nil {
		if errors.Is(err, usecase.ErrOrderAlreadyUploaded) {
			return c.String(http.StatusOK, "order has already been uploaded")
		}
		if errors.Is(err, usecase.ErrOrderAlreadyUploadedByAnotherUser) {
			return c.String(http.StatusConflict, "order has already been uploaded by another user")
		}
		s.logger.WithContext(ctx).WithError(err).Error("failed to create order")
		return err
	}

	return c.String(http.StatusAccepted, "order is uploaded")
}
