package restapi_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/georgg2003/gophermart/internal/delivery/restapi"
	"github.com/georgg2003/gophermart/internal/pkg/testutils"
	"github.com/georgg2003/gophermart/internal/usecase"
	"github.com/georgg2003/gophermart/pkg/gotils.go"
	"github.com/georgg2003/gophermart/pkg/jwthelper"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostAPIUserRegister(t *testing.T) {
	testApp := testutils.NewTestApp(t)

	helper := jwthelper.New([]byte(testApp.Cfg.JWTSecretKey))

	creds := restapi.RegisterRequest{
		Login:    "login",
		Password: "password",
	}
	credsBody, err := json.Marshal(creds)
	require.NoError(t, err)

	for _, tc := range []DeliveryTestCase{
		{
			name:       "success register",
			body:       credsBody,
			statusCode: http.StatusOK,
			response:   []byte("successfully registered"),
			mockFunc: func(req *http.Request) {
				testApp.Repo.EXPECT().CreateUser(
					req.Context(),
					creds.Login,
					gotils.HashPassword(creds.Password),
				).Return(testUserID, nil)
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
			name:       "no body",
			statusCode: http.StatusBadRequest,
			response:   []byte("wrong request format"),
		},
		{
			name:       "user already exists",
			body:       credsBody,
			statusCode: http.StatusConflict,
			response:   []byte("user already exists"),
			mockFunc: func(req *http.Request) {
				testApp.Repo.EXPECT().CreateUser(
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
				testApp.Repo.EXPECT().CreateUser(
					req.Context(),
					creds.Login,
					gotils.HashPassword(creds.Password),
				).Return(int64(-1), testutils.ErrUnexpectedError)
			},
			errExpected: true,
		},
	} {
		t.Run(tc.name, runDeliveryTestCase(tc, testApp.Server.PostAPIUserRegister))
	}
}
