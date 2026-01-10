package accrual

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/georgg2003/gophermart/internal/models"
	"github.com/georgg2003/gophermart/internal/pkg/config"
	"github.com/georgg2003/gophermart/internal/pkg/logging"
	"github.com/georgg2003/gophermart/internal/repository"
	"github.com/georgg2003/gophermart/pkg/errutils"
	"github.com/go-resty/resty/v2"
)

type accrual struct {
	client *resty.Client
	logger *logging.Logger
}

var ErrResponseNilBody = errors.New("accrual response has nil body")

func (a *accrual) GetOrderAccrual(
	ctx context.Context,
	orderNumber string,
) (*models.GetOrderAccrualResponse, error) {
	r, err := a.client.R().
		SetPathParam("number", orderNumber).
		Get("/api/orders/{number}")

	logger := a.logger.With(slog.String("order_number", orderNumber))
	logger.WithString("url", r.Request.URL).Debug("request to accrual service")

	if err != nil {
		logger.WithError(err).Error("request to accrual failed")
		return nil, errutils.Wrap(err, "request to accrual failed")
	}

	logger.WithString("response_body", r.String()).Debug("accrual response body")

	body := r.Body()
	if body == nil {
		logger.Error("accrual response has nil body")
		return nil, ErrResponseNilBody
	}

	var resp models.GetOrderAccrualResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		logger.WithError(err).Error("failed to unmarshall accrual response")
		return nil, errutils.Wrap(err, "failed to unmarshall accrual response")
	}

	return &resp, err
}

func New(cfg *config.Config, logger *logging.Logger) repository.AccrualRepo {
	client := resty.New().SetBaseURL(cfg.AccrualSysAddr).EnableTrace()
	return &accrual{
		client: client,
		logger: logger,
	}
}
