package restapi_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/georgg2003/gophermart/internal/delivery/restapi"
	"github.com/georgg2003/gophermart/internal/models"
	"github.com/georgg2003/gophermart/internal/pkg/testutils"
	"github.com/georgg2003/gophermart/internal/usecase"
	"github.com/georgg2003/gophermart/pkg/gotils.go"
	"github.com/georgg2003/gophermart/pkg/jwthelper"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostAPIUserLogin(t *testing.T) {
	testApp := testutils.NewTestApp(t)
	cfg := testApp.Cfg
	repo := testApp.Repo

	helper := jwthelper.New([]byte(cfg.JWTSecretKey))

	creds := restapi.LoginRequest{
		Login:    "login",
		Password: "password",
	}
	credsBody, err := json.Marshal(creds)
	require.NoError(t, err)

	for _, tc := range []DeliveryTestCase{
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
			validateResponse: func(t *testing.T, res *http.Response) {
				if res.StatusCode == http.StatusOK {
					token := res.Header.Get(echo.HeaderAuthorization)
					assert.NotEmpty(t, token)
					token = strings.Split(token, "Bearer ")[1]
					userID, err := helper.ReadAccessToken(token)
					assert.NoError(t, err)
					assert.Equal(t, userID, testUserID)
				}
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
				).Return(nil, testutils.UnexpectedError)
			},
			errExpected: true,
		},
		{
			name:       "invalid body",
			statusCode: http.StatusBadRequest,
			response:   []byte("wrong request format"),
		},
	} {
		t.Run(tc.name, runDeliveryTestCase(tc, testApp.Server.PostAPIUserLogin))
	}
}
