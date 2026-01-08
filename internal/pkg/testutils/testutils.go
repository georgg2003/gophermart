package testutils

import (
	"errors"
	"testing"

	"github.com/georgg2003/gophermart/internal/delivery/restapi"
	"github.com/georgg2003/gophermart/internal/pkg/config"
	"github.com/georgg2003/gophermart/internal/repository/mock"
	"github.com/georgg2003/gophermart/internal/usecase"
	"github.com/sirupsen/logrus"
	"go.uber.org/mock/gomock"
)

const TestOrderNumber = "12345678903"

var ErrUnexpectedError = errors.New("some error")

type TestApp struct {
	Cfg         *config.Config
	Repo        *mock.MockRepository
	AccrualRepo *mock.MockAccrualRepo
	Server      restapi.ServerInterface
	UseCase     usecase.UseCase
}

func NewTestApp(t *testing.T) *TestApp {
	cfg := config.New()
	logger := logrus.New()

	ctrl := gomock.NewController(t)
	repo := mock.NewMockRepository(ctrl)
	accrualRepo := mock.NewMockAccrualRepo(ctrl)

	uc := usecase.New(cfg, logger, repo, accrualRepo)
	server := restapi.NewServer(cfg, logger, uc)

	return &TestApp{
		Cfg:         cfg,
		Repo:        repo,
		AccrualRepo: accrualRepo,
		Server:      server,
		UseCase:     uc,
	}
}
