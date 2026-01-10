package usecase

import (
	"context"
	"time"

	"github.com/georgg2003/gophermart/internal/models"
	"github.com/sirupsen/logrus"
)

func (uc *useCase) processOrder(ctx context.Context) {
	orderNumber, err := uc.repo.GetOrderToProcess(ctx, uc.cfg.ProcessRetryTimeout)
	if err != nil {
		uc.logger.WithError(err).Error("failed to get order to process")
		return
	}
	if orderNumber == "" {
		uc.logger.Debug("no orders to process")
		return
	}

	logger := uc.logger.WithField("order_number", orderNumber)

	resp, err := uc.accrualRepo.GetOrderAccrual(ctx, orderNumber)
	if err != nil {
		logger.
			WithError(err).
			Error("failed to get order accrual")
		return
	}

	logger = logger.
		WithFields(logrus.Fields{
			"accrual_order_number": resp.Order,
			"accrual_order_status": resp.Status,
			"accrual_amount":       resp.Accrual,
		})

	switch resp.Status {
	case models.AccrualStatusProcessed:
		err = uc.repo.ApplyOrderAccrual(
			ctx,
			resp.Order,
			int(models.NewMoneyFromMajor(resp.Accrual).AmountMinor),
		)
	case models.AccrualStatusInvalid:
		err = uc.repo.SetOrderStatus(ctx, resp.Order, models.StatusInvalid)
	}
	if err != nil {
		logger.WithError(err).Error("failed to update order")
		return
	}
}

func (uc *useCase) workerIter(ctx context.Context) {
	newCtx, cancel := context.WithTimeout(ctx, time.Duration(uc.cfg.WorkerTimeout*int(time.Second)))
	defer cancel()

	uc.processOrder(newCtx)
}

func (uc *useCase) MakeProcessorWorker(
	ctx context.Context,
) {
	ticker := time.NewTicker(time.Duration(time.Second))
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			uc.logger.Info("shutting down processor worker")
			return
		case <-ticker.C:
			uc.workerIter(ctx)
		}
	}
}
