package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/georgg2003/gophermart/internal/models"
	"github.com/georgg2003/gophermart/internal/pkg/contextlib"
	"github.com/georgg2003/gophermart/pkg/errutils"
	"github.com/georgg2003/gophermart/pkg/gotils.go"
)

func (uc *useCase) UserRegister(
	ctx context.Context,
	login string,
	password string,
) (accessToken string, err error) {
	hash := gotils.HashPassword(password)
	logger := uc.logger.WithRequestCtx(ctx)
	userID, err := uc.repo.CreateUser(ctx, login, hash)
	if err != nil {
		return "", errutils.Wrap(err, "failed to create new user")
	}
	logger.Info("register success")

	return uc.jwtHelper.NewAccessToken(userID)
}

func (uc *useCase) UserLogin(
	ctx context.Context,
	login string,
	password string,
) (accessToken string, err error) {
	hash := gotils.HashPassword(password)
	logger := uc.logger.WithRequestCtx(ctx)
	userCredentials, err := uc.repo.GetUserByLogin(ctx, login)
	if err != nil {
		return "", errutils.Wrap(err, "failed to get user by login")
	}
	if userCredentials.PasswordHash != hash {
		return "", ErrUserWrongPassword
	}
	logger.Info("login success")

	return uc.jwtHelper.NewAccessToken(userCredentials.ID)
}

func (uc *useCase) UserGetBalance(
	ctx context.Context,
) (balance *models.UserBalance, err error) {
	userID := contextlib.MustGetUserID(ctx)
	logger := uc.logger.WithRequestCtx(ctx)
	balance, err = uc.repo.GetUserBalance(ctx, userID)
	if err != nil {
		return nil, errutils.Wrap(err, "failed to get user balance")
	}
	logger.With(
		slog.Float64("balance", balance.Current.Major()),
		slog.Float64("withdrawn", balance.Withdrawn.Major()),
	).Info("get balance success")

	return balance, err
}

func (uc *useCase) UserGetWithdrawals(
	ctx context.Context,
) (withdrawals []models.Withdrawal, err error) {
	userID := contextlib.MustGetUserID(ctx)
	logger := uc.logger.WithRequestCtx(ctx)
	withdrawals, err = uc.repo.GetUserWithdrawals(ctx, userID)
	if err != nil {
		return nil, errutils.Wrap(err, "failed to get user withdrawals")
	}
	if len(withdrawals) == 0 {
		return nil, ErrWidthdrawalsNotFound
	}
	logger.WithString("withdrawals", fmt.Sprint(withdrawals)).Info("get withdrawals success")

	return withdrawals, err
}

func (uc *useCase) UserGetOrders(
	ctx context.Context,
) (orders []models.Order, err error) {
	userID := contextlib.MustGetUserID(ctx)
	logger := uc.logger.WithRequestCtx(ctx)
	orders, err = uc.repo.GetUserOrders(ctx, userID)
	if err != nil {
		return nil, errutils.Wrap(err, "failed to get user orders")
	}
	if len(orders) == 0 {
		return nil, ErrOrdersNotFound
	}
	logger.WithString("orders", fmt.Sprint(orders)).Info("get orders success")

	return orders, err
}

func (uc *useCase) UserCreateOrder(
	ctx context.Context,
	orderNumber string,
) (err error) {
	userID := contextlib.MustGetUserID(ctx)
	logger := uc.logger.WithRequestCtx(ctx)
	err = uc.repo.CreateUserOrder(ctx, userID, orderNumber)
	if err != nil {
		err = errutils.Wrap(err, "failed to create order")
	}
	logger.WithString("order_number", orderNumber).Info("create order success")

	return err
}

func (uc *useCase) UserCreateWithdrawal(
	ctx context.Context,
	orderNumber string,
	amount models.Money,
) (err error) {
	userID := contextlib.MustGetUserID(ctx)
	logger := uc.logger.WithRequestCtx(ctx)
	err = uc.repo.CreateUserWithdrawal(ctx, userID, orderNumber, amount.AmountMinor)
	if err != nil {
		err = errutils.Wrap(err, "failed to withdraw")
	}
	logger.With(
		slog.Float64("amount", amount.Major()),
	).Info("success withdrawal")

	return err
}
