package repository

import (
	"context"

	"github.com/georgg2003/gophermart/internal/models"
)

//go:generate go tool mockgen -destination ./mock/mock.go -package mock . Repository,AccrualRepo
type Repository interface {
	CreateUser(
		ctx context.Context,
		login string,
		passwordHash string,
	) (id int64, err error)
	GetUserByLogin(
		ctx context.Context,
		login string,
	) (creds *models.UserCredentials, err error)
	GetUserBalance(
		ctx context.Context,
		userID int64,
	) (balance *models.UserBalance, err error)
	GetUserWithdrawals(
		ctx context.Context,
		userID int64,
	) (withdrawals []models.Withdrawal, err error)
	GetUserOrders(
		ctx context.Context,
		userID int64,
	) (orders []models.Order, err error)
	CreateUserOrder(
		ctx context.Context,
		userID int64,
		orderNumber string,
	) (err error)
	CreateUserWithdrawal(
		ctx context.Context,
		userID int64,
		orderNumber string,
		amount int64,
	) (err error)
	GetOrderToProcess(
		ctx context.Context,
		processRetryTimeout int,
	) (orderNumber string, err error)
	SetOrderStatus(
		ctx context.Context,
		orderNumber string,
		orderStatus models.OrderStatus,
	) (err error)
	ApplyOrderAccrual(
		ctx context.Context,
		orderNumber string,
		accrual int,
	) (err error)
}

type AccrualRepo interface {
	GetOrderAccrual(
		ctx context.Context,
		orderNumber string,
	) (response *models.GetOrderAccrualResponse, err error)
}
