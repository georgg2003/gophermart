package restapi_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

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

const getBalancePath = "/balance"

type GetBalanceTestCase struct {
	name        string
	body        []byte
	statusCode  int
	response    []byte
	mockFunc    func(*http.Request)
	errExpected bool
}

func runGetBalanceTestCase(
	tc GetBalanceTestCase,
	server restapi.ServerInterface,
) func(t *testing.T) {
	return func(t *testing.T) {
		buf := bytes.NewBuffer(tc.body)
		req := httptest.NewRequest(http.MethodGet, getBalancePath, buf)
		req = req.WithContext(contextlib.SetUserID(req.Context(), testUserID))

		resp := httptest.NewRecorder()

		c := echo.New().NewContext(req, resp)

		if tc.mockFunc != nil {
			tc.mockFunc(req)
		}

		err := server.GetAPIUserBalance(c)
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

func TestGetAPIUserBalance(t *testing.T) {
	cfg := config.New()

	logger := logrus.New()

	ctrl := gomock.NewController(t)
	repo := mock.NewMockRepository(ctrl)
	accrualRepo := mock.NewMockAccrualRepo(ctrl)

	uc := usecase.New(cfg, logger, repo, accrualRepo)
	server := restapi.NewServer(cfg, logger, uc)

	current := models.NewMoney(10000)
	withdrawn := models.NewMoney(1000)

	successResponse := restapi.BalanceResponse{
		Current:   current.Major(),
		Withdrawn: withdrawn.Major(),
	}
	response, err := json.Marshal(successResponse)
	response = append(response, '\n')
	require.NoError(t, err)

	for _, tc := range []GetBalanceTestCase{
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
				).Return(nil, errors.New("some error"))
			},
			errExpected: true,
		},
	} {
		t.Run(tc.name, runGetBalanceTestCase(tc, server))
	}
}
