package restapi_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/georgg2003/gophermart/internal/delivery/restapi"
	"github.com/georgg2003/gophermart/internal/models"
	"github.com/georgg2003/gophermart/internal/pkg/testutils"
	"github.com/georgg2003/gophermart/internal/usecase"
	"github.com/stretchr/testify/require"
)

func TestGetAPIUserWithdrawals(t *testing.T) {
	app := testutils.NewTestApp(t)
	repo := app.Repo

	orderNumber := "12345678903"
	processedAt := time.Now()
	money := models.NewMoney(100000)

	successResponse := []restapi.WithdrawalInfo{
		{
			Order:       &orderNumber,
			ProcessedAt: &processedAt,
			Sum:         money.Major(),
		},
	}
	response, err := json.Marshal(successResponse)
	response = append(response, '\n')
	require.NoError(t, err)

	for _, tc := range []DeliveryTestCase{
		{
			name:       "success get withdrawals",
			statusCode: http.StatusOK,
			response:   response,
			mockFunc: func(req *http.Request) {
				repo.EXPECT().GetUserWithdrawals(
					req.Context(),
					testUserID,
				).Return([]models.Withdrawal{
					{
						Order:       &orderNumber,
						Sum:         money,
						ProcessedAt: &processedAt,
					},
				}, nil)
			},
		},
		{
			name:       "no withdrawals",
			statusCode: http.StatusNoContent,
			response:   []byte("No withdrawals found"),
			mockFunc: func(req *http.Request) {
				repo.EXPECT().GetUserWithdrawals(
					req.Context(),
					testUserID,
				).Return(nil, usecase.ErrWidthdrawalsNotFound)
			},
		},
		{
			name: "repo error",
			mockFunc: func(req *http.Request) {
				repo.EXPECT().GetUserWithdrawals(
					req.Context(),
					testUserID,
				).Return(nil, testutils.UnexpectedError)
			},
			errExpected: true,
		},
	} {
		t.Run(tc.name, runDeliveryTestCase(tc, app.Server.GetAPIUserWithdrawals))
	}
}
