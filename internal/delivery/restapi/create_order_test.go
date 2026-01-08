package restapi_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/georgg2003/gophermart/internal/delivery/restapi"
	"github.com/georgg2003/gophermart/internal/usecase"
)

func TestPostAPIUserOrders(t *testing.T) {
	app := newTestApp(t)
	repo := app.repo

	req := restapi.PostAPIUserOrdersTextRequestBody(testOrderNumber)
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
					testOrderNumber,
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
					testOrderNumber,
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
					testOrderNumber,
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
					testOrderNumber,
				).Return(errors.New("some error"))
			},
			errExpected: true,
		},
	} {
		t.Run(tc.name, runDeliveryTestCase(tc, app.server.PostAPIUserOrders))
	}
}
