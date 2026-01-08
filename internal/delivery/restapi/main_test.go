package restapi_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/georgg2003/gophermart/internal/delivery/restapi"
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

const testUserID = int64(1)
const testOrderNumber = "123321"

type DeliveryTestCase struct {
	name             string
	body             []byte
	statusCode       int
	response         []byte
	mockFunc         func(*http.Request)
	errExpected      bool
	validateResponse func(t *testing.T, r *http.Response)
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

type testApp struct {
	cfg         *config.Config
	repo        *mock.MockRepository
	accrualRepo *mock.MockAccrualRepo
	server      restapi.ServerInterface
}

func newTestApp(t *testing.T) *testApp {
	cfg := config.New()
	logger := logrus.New()

	ctrl := gomock.NewController(t)
	repo := mock.NewMockRepository(ctrl)
	accrualRepo := mock.NewMockAccrualRepo(ctrl)

	uc := usecase.New(cfg, logger, repo, accrualRepo)
	server := restapi.NewServer(cfg, logger, uc)

	return &testApp{
		cfg:         cfg,
		repo:        repo,
		accrualRepo: accrualRepo,
		server:      server,
	}
}
