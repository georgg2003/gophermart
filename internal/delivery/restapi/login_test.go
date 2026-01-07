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
	"github.com/georgg2003/gophermart/internal/models"
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

const loginPath = "/login"

type LoginTestCase struct {
	name        string
	body        []byte
	statusCode  int
	response    []byte
	mockFunc    func(*http.Request)
	errExpected bool
}

func runLoginTestCase(
	tc LoginTestCase,
	server restapi.ServerInterface,
	helper *jwthelper.JWTHelper,
) func(t *testing.T) {
	return func(t *testing.T) {
		buf := bytes.NewBuffer(tc.body)
		req := httptest.NewRequest(http.MethodGet, loginPath, buf)

		resp := httptest.NewRecorder()

		c := echo.New().NewContext(req, resp)

		if tc.mockFunc != nil {
			tc.mockFunc(req)
		}

		err := server.PostAPIUserLogin(c)
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
		assert.Equal(t, tc.response, body)

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

func TestPostAPIUserLogin(t *testing.T) {
	cfg := config.New()

	logger := logrus.New()

	helper := jwthelper.New([]byte(cfg.JWTSecretKey))

	ctrl := gomock.NewController(t)
	repo := mock.NewMockRepository(ctrl)
	accrualRepo := mock.NewMockAccrualRepo(ctrl)

	uc := usecase.New(cfg, logger, repo, accrualRepo)
	server := restapi.NewServer(cfg, logger, uc)

	creds := restapi.LoginRequest{
		Login:    "login",
		Password: "password",
	}
	credsBody, err := json.Marshal(creds)
	require.NoError(t, err)

	for _, tc := range []LoginTestCase{
		{
			name:       "success login",
			body:       credsBody,
			statusCode: http.StatusOK,
			response:   []byte("successfully logged in"),
			mockFunc: func(req *http.Request) {
				repo.EXPECT().GetUserByLogin(
					req.Context(),
					creds.Login,
				).Return(&models.UserCredentials{
					ID:           testUserID,
					Login:        creds.Login,
					PasswordHash: gotils.HashPassword(creds.Password),
				}, nil)
			},
		},
		{
			name:       "wrong password",
			body:       credsBody,
			statusCode: http.StatusUnauthorized,
			response:   []byte("Incorrect login or password"),
			mockFunc: func(req *http.Request) {
				repo.EXPECT().GetUserByLogin(
					req.Context(),
					creds.Login,
				).Return(&models.UserCredentials{
					ID:           testUserID,
					Login:        creds.Login,
					PasswordHash: gotils.HashPassword("another password"),
				}, nil)
			},
		},
		{
			name:       "wrong user",
			body:       credsBody,
			statusCode: http.StatusUnauthorized,
			response:   []byte("Incorrect login or password"),
			mockFunc: func(req *http.Request) {
				repo.EXPECT().GetUserByLogin(
					req.Context(),
					creds.Login,
				).Return(nil, usecase.ErrUserNotFound)
			},
		},
		{
			name:       "repo error",
			body:       credsBody,
			statusCode: http.StatusInternalServerError,
			mockFunc: func(req *http.Request) {
				repo.EXPECT().GetUserByLogin(
					req.Context(),
					creds.Login,
				).Return(nil, errors.New("some error"))
			},
			errExpected: true,
		},
		{
			name:       "invalid body",
			statusCode: http.StatusBadRequest,
			response:   []byte("wrong request format"),
		},
	} {
		t.Run(tc.name, runLoginTestCase(tc, server, helper))
	}
}
