package restapi_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/georgg2003/gophermart/internal/delivery/restapi"
	"github.com/georgg2003/gophermart/internal/models"
	"github.com/georgg2003/gophermart/internal/pkg/config"
	"github.com/georgg2003/gophermart/internal/pkg/contextlib"
	"github.com/georgg2003/gophermart/internal/repository/mock"
	"github.com/georgg2003/gophermart/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const getOrdersPath = "/orders"

type GetOrdersTestCase struct {
	name        string
	body        []byte
	statusCode  int
	response    []byte
	mockFunc    func(*http.Request)
	errExpected bool
}

func runGetOrdersTestCase(
	tc GetOrdersTestCase,
	server restapi.ServerInterface,
) func(t *testing.T) {
	return func(t *testing.T) {
		buf := bytes.NewBuffer(tc.body)
		req := httptest.NewRequest(http.MethodGet, getOrdersPath, buf)
		req = req.WithContext(contextlib.SetUserID(req.Context(), testUserID))

		resp := httptest.NewRecorder()

		c := echo.New().NewContext(req, resp)

		if tc.mockFunc != nil {
			tc.mockFunc(req)
		}

		err := server.GetAPIUserOrders(c)
		if tc.errExpected {
			assert.Error(t, err)
			return
		}
		require.NoError(t, err)

		res := resp.Result()
		assert.Equal(t, tc.statusCode, res.StatusCode)

		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		require.NoError(t, err)
		assert.Equal(t, string(tc.response), string(body))
	}
}

func TestGetAPIUserOrders(t *testing.T) {
	cfg := config.New()

	logger := logrus.New()

	ctrl := gomock.NewController(t)
	repo := mock.NewMockRepository(ctrl)
	accrualRepo := mock.NewMockAccrualRepo(ctrl)

	uc := usecase.New(cfg, logger, repo, accrualRepo)
	server := restapi.NewServer(cfg, logger, uc)

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

	for _, tc := range []GetOrdersTestCase{
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
		t.Run(tc.name, runGetOrdersTestCase(tc, server))
	}
}
