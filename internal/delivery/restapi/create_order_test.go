package restapi_test

import (
	"net/http"
	"testing"

	"github.com/georgg2003/gophermart/internal/delivery/restapi"
	"github.com/georgg2003/gophermart/internal/pkg/testutils"
	"github.com/georgg2003/gophermart/internal/usecase"
)

func TestPostAPIUserOrders(t *testing.T) {
	app := testutils.NewTestApp(t)
	repo := app.Repo

	req := restapi.PostAPIUserOrdersTextRequestBody(testutils.TestOrderNumber)
	body := []byte(req)

	invalidReq := restapi.PostAPIUserOrdersTextRequestBody("321321321")
	invalidReqBody := []byte(invalidReq)

	for _, tc := range []DeliveryTestCase{
		{
			name:       "success create order",
			statusCode: http.StatusAccepted,
			body:       body,
			response:   []byte("order is uploaded"),
			mockFunc: func(req *http.Request) {
				repo.EXPECT().CreateUserOrder(
					req.Context(),
					testUserID,
					testutils.TestOrderNumber,
				).Return(nil)
			},
		},
		{
			name:       "invalid body",
			statusCode: http.StatusBadRequest,
			body:       []byte("long long long string"),
			transformRequest: func(r *http.Request, w http.ResponseWriter) {
				r.Body = http.MaxBytesReader(w, r.Body, 1)
			},
			response: []byte("wrong body format"),
		},
		{
			name:       "invalid order number",
			statusCode: http.StatusUnprocessableEntity,
			body:       invalidReqBody,
			response:   []byte("order number is not valid"),
		},
		{
			name:       "order already uploaded",
			body:       body,
			statusCode: http.StatusOK,
			response:   []byte("order has already been uploaded"),
			mockFunc: func(req *http.Request) {
				repo.EXPECT().CreateUserOrder(
					req.Context(),
					testUserID,
					testutils.TestOrderNumber,
				).Return(usecase.ErrOrderAlreadyUploaded)
			},
		},
		{
			name:       "order already uploaded by another user",
			body:       body,
			statusCode: http.StatusConflict,
			response:   []byte("order has already been uploaded by another user"),
			mockFunc: func(req *http.Request) {
				repo.EXPECT().CreateUserOrder(
					req.Context(),
					testUserID,
					testutils.TestOrderNumber,
				).Return(usecase.ErrOrderAlreadyUploadedByAnotherUser)
			},
		},
		{
			name: "repo error",
			body: body,
			mockFunc: func(req *http.Request) {
				repo.EXPECT().CreateUserOrder(
					req.Context(),
					testUserID,
					testutils.TestOrderNumber,
				).Return(testutils.UnexpectedError)
			},
			errExpected: true,
		},
	} {
		t.Run(tc.name, runDeliveryTestCase(tc, app.Server.PostAPIUserOrders))
	}
}
