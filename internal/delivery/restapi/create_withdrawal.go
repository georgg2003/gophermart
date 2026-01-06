package restapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/georgg2003/gophermart/internal/models"
	"github.com/georgg2003/gophermart/internal/usecase"
	"github.com/georgg2003/gophermart/pkg/luhn"
	"github.com/labstack/echo/v4"
)

func (s *server) PostAPIUserBalanceWithdraw(c echo.Context) error {
	req := c.Request()
	defer req.Body.Close()
	ctx := req.Context()

	decoder := json.NewDecoder(req.Body)

	var withdrawRequest WithdrawRequest
	err := decoder.Decode(&withdrawRequest)
	if err != nil {
		s.logger.WithError(err).Info("failed to decode json withdraw request")
		return c.String(http.StatusBadRequest, "wrong request format")
	}

	if !luhn.ValidLuhn(withdrawRequest.Order) {
		s.logger.WithError(err).Info("order number is not valid")
		return c.String(http.StatusUnprocessableEntity, "order number is not valid")
	}

	err = s.uc.UserCreateWithdrawal(
		ctx,
		withdrawRequest.Order,
		models.NewMoneyFromMajor(withdrawRequest.Sum),
	)
	if err != nil {
		if errors.Is(err, usecase.ErrNotEnoughBalance) {
			s.logger.WithError(err).Info("not enough balance")
			return c.String(http.StatusPaymentRequired, "not enough balance for withdrawal")
		}
		if errors.Is(err, usecase.ErrWithdrawalAlreadyExists) {
			s.logger.WithError(err).Info("withdrawal already exists")
			return c.String(http.StatusConflict, "withdrawal already exists")
		}
		return err
	}
	return nil
}
