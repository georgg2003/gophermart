package accrual

import (
	"context"
	"encoding/json"

	"github.com/georgg2003/gophermart/internal/models"
	"github.com/georgg2003/gophermart/internal/pkg/config"
	"github.com/georgg2003/gophermart/internal/repository"
	"github.com/georgg2003/gophermart/pkg/errutils"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

type accrual struct {
	client *resty.Client
	logger *logrus.Logger
}

func (a *accrual) GetOrderAccrual(
	ctx context.Context,
	orderNumber string,
) (resp *models.GetOrderAccrualResponse, err error) {
	r, err := a.client.R().
		SetPathParam("number", orderNumber).
		Get("/api/orders/{number}")
	logger := a.logger.WithField("order_number", orderNumber)
	logger.WithField("trace", r.Request.TraceInfo()).Debug("request to accrual service")
	if err != nil {
		logger.WithError(err).Error("request to accrual failed")
		return nil, errutils.Wrap(err, "request to accrual failed")
	}
	body := r.Body()
	if body == nil {
		logger.Warn("accrual response has nil body")
		return nil, nil
	}
	if err = json.Unmarshal(r.Body(), resp); err != nil {
		logger.WithError(err).Error("failed to unmarshall accrual response")
		return nil, errutils.Wrap(err, "failed to unmarshall accrual response")
	}

	return resp, err
}

func New(cfg *config.Config, logger *logrus.Logger) repository.AccrualRepo {
	client := resty.New().SetBaseURL(cfg.AccrualSysAddr).EnableTrace()
	return &accrual{
		client: client,
		logger: logger,
	}
}
