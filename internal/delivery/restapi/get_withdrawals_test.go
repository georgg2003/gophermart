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
	"github.com/georgg2003/gophermart/pkg/jwthelper"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const getWithdrawalsPath = "/withdrawals"

type GetWithdrawalsTestCase struct {
	name        string
	body        []byte
	statusCode  int
	response    []byte
	mockFunc    func(*http.Request)
	errExpected bool
}

func runGetWithdrawalsTestCase(
	tc GetWithdrawalsTestCase,
	server restapi.ServerInterface,
	helper *jwthelper.JWTHelper,
) func(t *testing.T) {
	return func(t *testing.T) {
		buf := bytes.NewBuffer(tc.body)
		req := httptest.NewRequest(http.MethodGet, loginPath, buf)
		req = req.WithContext(contextlib.SetUserID(req.Context(), testUserID))

		resp := httptest.NewRecorder()

		c := echo.New().NewContext(req, resp)

		if tc.mockFunc != nil {
			tc.mockFunc(req)
		}

		err := server.GetAPIUserWithdrawals(c)
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

func TestGetAPIUserWithdrawals(t *testing.T) {
	cfg := config.New()

	logger := logrus.New()

	helper := jwthelper.New([]byte(cfg.JWTSecretKey))

	ctrl := gomock.NewController(t)
	repo := mock.NewMockRepository(ctrl)
	accrualRepo := mock.NewMockAccrualRepo(ctrl)

	uc := usecase.New(cfg, logger, repo, accrualRepo)
	server := restapi.NewServer(cfg, logger, uc)

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

	for _, tc := range []GetWithdrawalsTestCase{
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
				).Return(nil, errors.New("some error"))
			},
			errExpected: true,
		},
	} {
		t.Run(tc.name, runGetWithdrawalsTestCase(tc, server, helper))
	}
}
