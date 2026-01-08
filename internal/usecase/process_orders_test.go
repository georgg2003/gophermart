package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/georgg2003/gophermart/internal/models"
	"github.com/georgg2003/gophermart/internal/pkg/testutils"
	"go.uber.org/mock/gomock"
)

type testCase struct {
	name     string
	mockFunc func()
}

func runWorkerTestCase(app *testutils.TestApp, tc testCase) func(t *testing.T) {
	return func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		done := make(chan struct{})

		if tc.mockFunc != nil {
			tc.mockFunc()
		}

		go func() {
			app.UseCase.MakeProcessorWorker(ctx)
			close(done)
		}()

		time.Sleep(1100 * time.Millisecond)

		cancel()

		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Fatal("worker did not stop after context cancellation")
		}
	}
}

func TestMakeProccessorWorker(t *testing.T) {
	app := testutils.NewTestApp(t)

	for _, tc := range []testCase{
		{
			name: "successful proccessing",
			mockFunc: func() {
				getOrderCall := app.Repo.EXPECT().
					GetOrderToProcess(gomock.Any(), app.Cfg.ProcessRetryTimeout).
					Return(testutils.TestOrderNumber, nil)

				money := models.NewMoney(10000)

				accrualResponse := &models.GetOrderAccrualResponse{
					Order:   testutils.TestOrderNumber,
					Status:  models.AccrualStatusProcessed,
					Accrual: money.Major(),
				}

				getOrderAccrualCall := app.AccrualRepo.EXPECT().
					GetOrderAccrual(gomock.Any(), testutils.TestOrderNumber).
					After(getOrderCall).
					Return(accrualResponse, nil)

				app.Repo.EXPECT().
					ApplyOrderAccrual(gomock.Any(), testutils.TestOrderNumber, int(money.AmountMinor)).
					After(getOrderAccrualCall).
					Return(nil)
			},
		},
		{
			name: "accrual status invalid",
			mockFunc: func() {
				getOrderCall := app.Repo.EXPECT().
					GetOrderToProcess(gomock.Any(), app.Cfg.ProcessRetryTimeout).
					Return(testutils.TestOrderNumber, nil)

				accrualResponse := &models.GetOrderAccrualResponse{
					Order:  testutils.TestOrderNumber,
					Status: models.AccrualStatusInvalid,
				}

				getOrderAccrualCall := app.AccrualRepo.EXPECT().
					GetOrderAccrual(gomock.Any(), testutils.TestOrderNumber).
					After(getOrderCall).
					Return(accrualResponse, nil)

				app.Repo.EXPECT().
					SetOrderStatus(gomock.Any(), testutils.TestOrderNumber, models.StatusInvalid).
					After(getOrderAccrualCall).
					Return(nil)
			},
		},
		{
			name: "failed to update order",
			mockFunc: func() {
				getOrderCall := app.Repo.EXPECT().
					GetOrderToProcess(gomock.Any(), app.Cfg.ProcessRetryTimeout).
					Return(testutils.TestOrderNumber, nil)

				accrualResponse := &models.GetOrderAccrualResponse{
					Order:  testutils.TestOrderNumber,
					Status: models.AccrualStatusInvalid,
				}

				getOrderAccrualCall := app.AccrualRepo.EXPECT().
					GetOrderAccrual(gomock.Any(), testutils.TestOrderNumber).
					After(getOrderCall).
					Return(accrualResponse, nil)

				app.Repo.EXPECT().
					SetOrderStatus(gomock.Any(), testutils.TestOrderNumber, models.StatusInvalid).
					After(getOrderAccrualCall).
					Return(testutils.ErrUnexpectedError)
			},
		},
		{
			name: "no orders to process",
			mockFunc: func() {
				app.Repo.EXPECT().
					GetOrderToProcess(gomock.Any(), app.Cfg.ProcessRetryTimeout).
					Return("", nil)

				app.AccrualRepo.EXPECT().
					GetOrderAccrual(gomock.Any(), testutils.TestOrderNumber).Times(0)

				app.Repo.EXPECT().
					ApplyOrderAccrual(gomock.Any(), testutils.TestOrderNumber, gomock.Any()).Times(0)
			},
		},
		{
			name: "can't get order to process",
			mockFunc: func() {
				app.Repo.EXPECT().
					GetOrderToProcess(gomock.Any(), app.Cfg.ProcessRetryTimeout).
					Return("", testutils.ErrUnexpectedError)

				app.AccrualRepo.EXPECT().
					GetOrderAccrual(gomock.Any(), testutils.TestOrderNumber).Times(0)

				app.Repo.EXPECT().
					ApplyOrderAccrual(gomock.Any(), testutils.TestOrderNumber, gomock.Any()).Times(0)
			},
		},
		{
			name: "failed to get accrual",
			mockFunc: func() {
				getOrderCall := app.Repo.EXPECT().
					GetOrderToProcess(gomock.Any(), app.Cfg.ProcessRetryTimeout).
					Return(testutils.TestOrderNumber, nil)

				app.AccrualRepo.EXPECT().
					GetOrderAccrual(gomock.Any(), testutils.TestOrderNumber).
					After(getOrderCall).
					Return(nil, testutils.ErrUnexpectedError)

				app.Repo.EXPECT().
					ApplyOrderAccrual(gomock.Any(), testutils.TestOrderNumber, gomock.Any()).
					Times(0)
			},
		},
	} {
		t.Run(tc.name, runWorkerTestCase(app, tc))
	}
}
