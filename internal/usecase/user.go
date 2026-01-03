package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/georgg2003/gophermart/internal/models"
	"github.com/georgg2003/gophermart/internal/pkg/contextlib"
	"github.com/georgg2003/gophermart/pkg/errutils"
)

func passwordToHash(password string) string {
	hash := sha256.New()
	return hex.EncodeToString(hash.Sum([]byte(password)))
}

func (uc *useCase) UserRegister(
	ctx context.Context,
	login string,
	password string,
) (accessToken string, err error) {
	hash := passwordToHash(password)
	userID, err := uc.repo.CreateUser(ctx, login, hash)
	if err != nil {
		return "", errutils.Wrap(err, "failed to create new user")
	}

	return uc.jwtHelper.NewAccessToken(userID)
}

func (uc *useCase) UserLogin(
	ctx context.Context,
	login string,
	password string,
) (accessToken string, err error) {
	hash := passwordToHash(password)
	userCredentials, err := uc.repo.GetUserByLogin(ctx, login)
	if err != nil {
		return "", errutils.Wrap(err, "failed to get user by login")
	}
	if userCredentials.PasswordHash != hash {
		return "", ErrUserWrongPassword
	}

	return uc.jwtHelper.NewAccessToken(userCredentials.ID)
}

func (uc *useCase) UserGetBalance(
	ctx context.Context,
) (balance *models.UserBalance, err error) {
	userID := contextlib.MustGetUserID(ctx, uc.logger)
	balance, err = uc.repo.GetUserBalance(ctx, userID)
	if err != nil {
		return nil, errutils.Wrap(err, "failed to get user balance")
	}
	return balance, err
}

func (uc *useCase) UserGetWithdrawals(
	ctx context.Context,
) (withdrawals []models.Withdrawal, err error) {
	userID := contextlib.MustGetUserID(ctx, uc.logger)
	withdrawals, err = uc.repo.GetUserWithdrawals(ctx, userID)
	if err != nil {
		return nil, errutils.Wrap(err, "failed to get user withdrawals")
	}
	if len(withdrawals) == 0 {
		return nil, ErrWidthdrawalsNotFound
	}
	return withdrawals, err
}

func (uc *useCase) UserGetOrders(
	ctx context.Context,
) (orders []models.Order, err error) {
	userID := contextlib.MustGetUserID(ctx, uc.logger)
	orders, err = uc.repo.GetUserOrders(ctx, userID)
	if err != nil {
		return nil, errutils.Wrap(err, "failed to get user withdrawals")
	}
	if len(orders) == 0 {
		return nil, ErrOrdersNotFound
	}
	return orders, err
}

func (uc *useCase) UserCreateOrder(
	ctx context.Context,
	orderNumber string,
) (err error) {
	userID := contextlib.MustGetUserID(ctx, uc.logger)
	err = uc.repo.CreateUserOrder(ctx, userID, orderNumber)
	if err != nil {
		err = errutils.Wrap(err, "failed to create order")
	}
	return err
}

func (uc *useCase) UserCreateWithdrawal(
	ctx context.Context,
	orderNumber string,
	amount models.Money,
) (err error) {
	userID := contextlib.MustGetUserID(ctx, uc.logger)
	err = uc.repo.CreateUserWithdrawal(ctx, userID, orderNumber, amount.AmountMinor)
	if err != nil {
		err = errutils.Wrap(err, "failed to withdraw")
	}
	return err
}
