package middleware_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/georgg2003/gophermart/internal/pkg/contextlib"
	"github.com/georgg2003/gophermart/internal/pkg/middleware"
	"github.com/georgg2003/gophermart/pkg/jwthelper"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const secretKey = "secret"
const testUserID = int64(1)
const skipperPath = "/register"
const successAuthTestCaseName = "success auth"
const invalidTokenResponse = "Invalid token"

type TestCase struct {
	name          string
	authorization string
	path          string
	statusCode    int
	response      []byte
}

func runTestCase(tc TestCase, handlerFunc echo.HandlerFunc) func(t *testing.T) {
	return func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, tc.path, nil)
		if tc.authorization != "" {
			req.Header.Set(echo.HeaderAuthorization, tc.authorization)
		}

		resp := httptest.NewRecorder()

		c := echo.New().NewContext(req, resp)
		c.SetPath(tc.path)

		err := handlerFunc(c)
		require.NoError(t, err)

		res := resp.Result()
		assert.Equal(t, tc.statusCode, res.StatusCode)

		if tc.name == successAuthTestCaseName {
			ctx := c.Request().Context()
			userID, ok := contextlib.GetUserID(ctx)
			assert.True(t, ok)
			assert.Equal(t, testUserID, userID)
		}

		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		assert.Equal(t, tc.response, body)
	}
}

func TestNewAuthMiddleware(t *testing.T) {
	logger := logrus.New()

	helper := jwthelper.New([]byte(secretKey))

	authMiddleware := middleware.NewAuthMiddleware(secretKey, logger, func(c echo.Context) bool {
		if c.Path() == skipperPath {
			return true
		}
		return false
	})

	next := func(c echo.Context) error {
		return nil
	}

	handlerFunc := authMiddleware(next)

	token, err := helper.NewAccessToken(testUserID)
	require.NoError(t, err)

	for _, tc := range []TestCase{
		{
			name:          successAuthTestCaseName,
			authorization: fmt.Sprintf("Bearer %s", token),
			path:          "/test",
			statusCode:    http.StatusOK,
			response:      []byte{},
		},
		{
			name:          "no auth",
			authorization: "",
			path:          "/test",
			statusCode:    http.StatusUnauthorized,
			response:      []byte("No auth"),
		},
		{
			name:          "skip auth",
			authorization: "invalid",
			path:          skipperPath,
			statusCode:    http.StatusOK,
			response:      []byte{},
		},
		{
			name:          "not bearer",
			authorization: token,
			path:          "/test",
			statusCode:    http.StatusUnauthorized,
			response:      []byte(invalidTokenResponse),
		},
		{
			name:          "invalid token",
			authorization: "invalid",
			path:          "/test",
			statusCode:    http.StatusUnauthorized,
			response:      []byte(invalidTokenResponse),
		},
	} {
		t.Run(tc.name, runTestCase(tc, handlerFunc))
	}
}
