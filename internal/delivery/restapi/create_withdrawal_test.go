package restapi_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/georgg2003/gophermart/internal/delivery/restapi"
	"github.com/georgg2003/gophermart/internal/models"
	"github.com/georgg2003/gophermart/internal/pkg/testutils"
	"github.com/georgg2003/gophermart/internal/usecase"
	"github.com/stretchr/testify/require"
)

func TestPostAPIUserBalanceWithdraw(t *testing.T) {
	app := testutils.NewTestApp(t)
	repo := app.Repo

	withdrawMoney := models.NewMoney(10000)

	req := restapi.WithdrawRequest{
		Order: testutils.TestOrderNumber,
		Sum:   withdrawMoney.Major(),
	}
	reqBody, err := json.Marshal(req)
	require.NoError(t, err)

	invalidReq := restapi.WithdrawRequest{
		Order: "213321",
		Sum:   withdrawMoney.Major(),
	}
	invalidReqBody, err := json.Marshal(invalidReq)
	require.NoError(t, err)

	for _, tc := range []DeliveryTestCase{
		{
			name:       "success create withdrawal",
			statusCode: http.StatusOK,
			body:       reqBody,
			mockFunc: func(req *http.Request) {
				repo.EXPECT().CreateUserWithdrawal(
					req.Context(),
					testUserID,
					testutils.TestOrderNumber,
					withdrawMoney.AmountMinor,
				).Return(nil)
			},
		},
		{
			name:       "invalid body",
			statusCode: http.StatusBadRequest,
			response:   []byte("wrong request format"),
		},
		{
			name:       "invalid order number",
			statusCode: http.StatusUnprocessableEntity,
			body:       invalidReqBody,
			response:   []byte("order number is not valid"),
		},
		{
			name:       "not enough balance",
			body:       reqBody,
			statusCode: http.StatusPaymentRequired,
			response:   []byte("not enough balance for withdrawal"),
			mockFunc: func(req *http.Request) {
				repo.EXPECT().CreateUserWithdrawal(
					req.Context(),
					testUserID,
					testutils.TestOrderNumber,
					withdrawMoney.AmountMinor,
				).Return(usecase.ErrNotEnoughBalance)
			},
		},
		{
			name:       "withdrawal exists",
			body:       reqBody,
			statusCode: http.StatusConflict,
			response:   []byte("withdrawal already exists"),
			mockFunc: func(req *http.Request) {
				repo.EXPECT().CreateUserWithdrawal(
					req.Context(),
					testUserID,
					testutils.TestOrderNumber,
					withdrawMoney.AmountMinor,
				).Return(usecase.ErrWithdrawalAlreadyExists)
			},
		},
		{
			name: "repo error",
			body: reqBody,
			mockFunc: func(req *http.Request) {
				repo.EXPECT().CreateUserWithdrawal(
					req.Context(),
					testUserID,
					testutils.TestOrderNumber,
					withdrawMoney.AmountMinor,
				).Return(testutils.UnexpectedError)
			},
			errExpected: true,
		},
	} {
		t.Run(tc.name, runDeliveryTestCase(tc, app.Server.PostAPIUserBalanceWithdraw))
	}
}
