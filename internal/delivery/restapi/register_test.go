package restapi_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/georgg2003/gophermart/internal/delivery/restapi"
	"github.com/georgg2003/gophermart/internal/pkg/config"
	"github.com/georgg2003/gophermart/internal/repository/mock"
	"github.com/georgg2003/gophermart/internal/usecase"
	"github.com/georgg2003/gophermart/pkg/gotils.go"
	"github.com/georgg2003/gophermart/pkg/jwthelper"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const registerPath = "/register"

type RegisterTestCase struct {
	name        string
	body        []byte
	statusCode  int
	response    []byte
	mockFunc    func(*http.Request)
	errExpected bool
}

func runRegisterTestCase(
	tc RegisterTestCase,
	server restapi.ServerInterface,
	helper *jwthelper.JWTHelper,
) func(t *testing.T) {
	return func(t *testing.T) {
		buf := bytes.NewBuffer(tc.body)
		req := httptest.NewRequest(http.MethodGet, registerPath, buf)

		resp := httptest.NewRecorder()

		c := echo.New().NewContext(req, resp)

		if tc.mockFunc != nil {
			tc.mockFunc(req)
		}

		err := server.PostAPIUserRegister(c)
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

		if res.StatusCode == http.StatusOK {
			token := res.Header.Get(echo.HeaderAuthorization)
			assert.NotEmpty(t, token)
			token = strings.Split(token, "Bearer ")[1]
			userID, err := helper.ReadAccessToken(token)
			assert.NoError(t, err)
			assert.Equal(t, userID, testUserID)
		}
	}
}

func TestPostAPIUserRegister(t *testing.T) {
	cfg := config.New()

	logger := logrus.New()

	helper := jwthelper.New([]byte(cfg.JWTSecretKey))

	ctrl := gomock.NewController(t)
	repo := mock.NewMockRepository(ctrl)
	accrualRepo := mock.NewMockAccrualRepo(ctrl)

	uc := usecase.New(cfg, logger, repo, accrualRepo)
	server := restapi.NewServer(cfg, logger, uc)

	creds := restapi.RegisterRequest{
		Login:    "login",
		Password: "password",
	}
	credsBody, err := json.Marshal(creds)
	require.NoError(t, err)

	for _, tc := range []RegisterTestCase{
		{
			name:       "success register",
			body:       credsBody,
			statusCode: http.StatusOK,
			response:   []byte("successfully registered"),
			mockFunc: func(req *http.Request) {
				repo.EXPECT().CreateUser(
					req.Context(),
					creds.Login,
					gotils.HashPassword(creds.Password),
				).Return(testUserID, nil)
			},
		},
		{
			name:       "success register",
			statusCode: http.StatusBadRequest,
			response:   []byte("wrong request format"),
		},
		{
			name:       "user already exists",
			body:       credsBody,
			statusCode: http.StatusConflict,
			response:   []byte("user already exists"),
			mockFunc: func(req *http.Request) {
				repo.EXPECT().CreateUser(
					req.Context(),
					creds.Login,
					gotils.HashPassword(creds.Password),
				).Return(int64(-1), usecase.ErrUserAlreadyExists)
			},
		},
		{
			name:       "insert fail",
			body:       credsBody,
			statusCode: http.StatusInternalServerError,
			mockFunc: func(req *http.Request) {
				repo.EXPECT().CreateUser(
					req.Context(),
					creds.Login,
					gotils.HashPassword(creds.Password),
				).Return(int64(-1), errors.New("some error"))
			},
			errExpected: true,
		},
	} {
		t.Run(tc.name, runRegisterTestCase(tc, server, helper))
	}
}
