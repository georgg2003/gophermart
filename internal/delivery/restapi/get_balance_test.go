package restapi_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/georgg2003/gophermart/internal/delivery/restapi"
	"github.com/georgg2003/gophermart/internal/models"
	"github.com/georgg2003/gophermart/internal/pkg/testutils"
	"github.com/stretchr/testify/require"
)

func TestGetAPIUserBalance(t *testing.T) {
	app := testutils.NewTestApp(t)
	repo := app.Repo

	current := models.NewMoney(10000)
	withdrawn := models.NewMoney(1000)

	successResponse := restapi.BalanceResponse{
		Current:   current.Major(),
		Withdrawn: withdrawn.Major(),
	}
	response, err := json.Marshal(successResponse)
	response = append(response, '\n')
	require.NoError(t, err)

	for _, tc := range []DeliveryTestCase{
		{
			name:       "success get balance",
			statusCode: http.StatusOK,
			response:   response,
			mockFunc: func(req *http.Request) {
				repo.EXPECT().GetUserBalance(
					req.Context(),
					testUserID,
				).Return(&models.UserBalance{
					Current:   current,
					Withdrawn: withdrawn,
				}, nil)
			},
		},
		{
			name: "repo error",
			mockFunc: func(req *http.Request) {
				repo.EXPECT().GetUserBalance(
					req.Context(),
					testUserID,
				).Return(nil, testutils.UnexpectedError)
			},
			errExpected: true,
		},
	} {
		t.Run(tc.name, runDeliveryTestCase(tc, app.Server.GetAPIUserBalance))
	}
}
