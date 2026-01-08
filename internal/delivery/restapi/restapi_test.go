package restapi_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/georgg2003/gophermart/internal/pkg/contextlib"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testUserID = int64(1)

type DeliveryTestCase struct {
	name             string
	body             []byte
	statusCode       int
	response         []byte
	mockFunc         func(*http.Request)
	errExpected      bool
	validateResponse func(t *testing.T, r *http.Response)
	transformRequest func(r *http.Request, w http.ResponseWriter)
}

func runDeliveryTestCase(
	tc DeliveryTestCase,
	callMethod func(c echo.Context) error,
) func(t *testing.T) {
	return func(t *testing.T) {
		buf := bytes.NewBuffer(tc.body)
		req := httptest.NewRequest(http.MethodGet, "/test", buf)
		req = req.WithContext(contextlib.SetUserID(req.Context(), testUserID))

		resp := httptest.NewRecorder()

		if tc.transformRequest != nil {
			tc.transformRequest(req, resp)
		}

		c := echo.New().NewContext(req, resp)

		if tc.mockFunc != nil {
			tc.mockFunc(req)
		}

		err := callMethod(c)
		if tc.errExpected {
			assert.Error(t, err)
			return
		}
		require.NoError(t, err)

		res := resp.Result()
		defer res.Body.Close()

		assert.Equal(t, tc.statusCode, res.StatusCode)

		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		require.NoError(t, err)
		assert.Equal(t, string(tc.response), string(body))

		if tc.validateResponse != nil {
			tc.validateResponse(t, res)
		}
	}
}
