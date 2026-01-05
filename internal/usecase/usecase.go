package usecase

import (
	"context"

	"github.com/georgg2003/gophermart/internal/models"
	"github.com/georgg2003/gophermart/internal/pkg/config"
	"github.com/georgg2003/gophermart/internal/repository"
	"github.com/georgg2003/gophermart/pkg/jwthelper"
	"github.com/sirupsen/logrus"
)

type useCase struct {
	cfg         *config.Config
	logger      *logrus.Logger
	repo        repository.Repository
	jwtHelper   *jwthelper.JWTHelper
	accrualRepo repository.AccrualRepo
}

type UseCase interface {
	UserRegister(
		ctx context.Context,
		login string,
		password string,
	) (accessToken string, err error)
	UserLogin(
		ctx context.Context,
		login string,
		password string,
	) (accessToken string, err error)
	UserGetBalance(
		ctx context.Context,
	) (balance *models.UserBalance, err error)
	UserGetWithdrawals(
		ctx context.Context,
	) (withdrawals []models.Withdrawal, err error)
	UserGetOrders(
		ctx context.Context,
	) (orders []models.Order, err error)
	UserCreateOrder(
		ctx context.Context,
		orderNumber string,
	) (err error)
	UserCreateWithdrawal(
		ctx context.Context,
		orderNumber string,
		amount models.Money,
	) (err error)
	MakeProcessorWorker(
		ctx context.Context,
	)
}

func New(
	cfg *config.Config,
	logger *logrus.Logger,
	repo repository.Repository,
	accrualRepo repository.AccrualRepo,
) UseCase {
	return &useCase{
		cfg:         cfg,
		logger:      logger,
		repo:        repo,
		jwtHelper:   jwthelper.New([]byte(cfg.JWTSecretKey)),
		accrualRepo: accrualRepo,
	}
}
