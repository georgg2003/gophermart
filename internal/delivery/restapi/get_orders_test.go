package restapi_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/georgg2003/gophermart/internal/delivery/restapi"
	"github.com/georgg2003/gophermart/internal/models"
	"github.com/georgg2003/gophermart/internal/usecase"
	"github.com/stretchr/testify/require"
)

func TestGetAPIUserOrders(t *testing.T) {
	app := newTestApp(t)
	repo := app.repo

	orderNumber := "12345678903"
	uploadedAt := time.Now().Add(-time.Hour)
	money := models.NewMoney(100000)

	successResponse := []restapi.OrderInfo{
		{
			Number:     orderNumber,
			Status:     restapi.OrderInfoStatus(models.StatusProcessed),
			UploadedAt: uploadedAt,
			Accrual:    money.NullMajor(),
		},
	}
	response, err := json.Marshal(successResponse)
	response = append(response, '\n')
	require.NoError(t, err)

	for _, tc := range []DeliveryTestCase{
		{
			name:       "success get orders",
			statusCode: http.StatusOK,
			response:   response,
			mockFunc: func(req *http.Request) {
				repo.EXPECT().GetUserOrders(
					req.Context(),
					testUserID,
				).Return([]models.Order{
					{
						Number:     orderNumber,
						Status:     models.StatusProcessed,
						Accrual:    &money,
						UploadedAt: uploadedAt,
					},
				}, nil)
			},
		},
		{
			name:       "no orders",
			statusCode: http.StatusNoContent,
			response:   []byte("No orders found"),
			mockFunc: func(req *http.Request) {
				repo.EXPECT().GetUserOrders(
					req.Context(),
					testUserID,
				).Return(nil, usecase.ErrOrdersNotFound)
			},
		},
		{
			name: "repo error",
			mockFunc: func(req *http.Request) {
				repo.EXPECT().GetUserOrders(
					req.Context(),
					testUserID,
				).Return(nil, errors.New("some error"))
			},
			errExpected: true,
		},
	} {
		t.Run(tc.name, runDeliveryTestCase(tc, app.server.GetAPIUserOrders))
	}
}
