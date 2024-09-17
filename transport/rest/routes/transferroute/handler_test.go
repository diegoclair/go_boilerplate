package transferroute_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/diegoclair/go_boilerplate/application/dto"
	"github.com/diegoclair/go_boilerplate/domain/entity"
	"github.com/diegoclair/go_boilerplate/transport/rest/routes/test"
	"github.com/diegoclair/go_boilerplate/transport/rest/routes/transferroute"
	"github.com/diegoclair/go_boilerplate/transport/rest/viewmodel"
	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/twinj/uuid"
)

func TestHandler_handleAddTransfer(t *testing.T) {
	body := viewmodel.TransferReq{
		AccountDestinationUUID: "randomUUID",
		Amount:                 5.55,
	}

	tests := append(test.PrivateEndpointValidations,
		test.PrivateEndpointTest{
			Name: "Should complete request with no error",
			Body: body,
			SetupAuth: func(ctx context.Context, t *testing.T, req *http.Request, m test.SvcMocks) {
				test.AddAuthorization(ctx, t, req, m)
			},
			BuildMocks: func(ctx context.Context, m test.SvcMocks, body any) {
				b := body.(viewmodel.TransferReq)
				m.TransferAppMock.EXPECT().CreateTransfer(ctx,
					dto.TransferInput{AccountDestinationUUID: b.AccountDestinationUUID, Amount: b.Amount}).
					Return(nil).MinTimes(1)
			},
			CheckResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, resp.Code)
				require.Empty(t, resp.Body)
			},
		},
		test.PrivateEndpointTest{
			Name: "Should return error if body is invalid",
			Body: "invalid body",
			SetupAuth: func(ctx context.Context, t *testing.T, req *http.Request, m test.SvcMocks) {
				test.AddAuthorization(ctx, t, req, m)
			},
			CheckResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, resp.Code)
				require.Contains(t, resp.Body.String(), "Unmarshal type error")
			},
		},
		test.PrivateEndpointTest{
			Name: "Should return error id we have some error on create transfer",
			Body: body,
			SetupAuth: func(ctx context.Context, t *testing.T, req *http.Request, m test.SvcMocks) {
				test.AddAuthorization(ctx, t, req, m)
			},
			BuildMocks: func(ctx context.Context, m test.SvcMocks, body any) {
				b := body.(viewmodel.TransferReq)
				m.TransferAppMock.EXPECT().CreateTransfer(ctx,
					dto.TransferInput{AccountDestinationUUID: b.AccountDestinationUUID, Amount: b.Amount}).
					Return(fmt.Errorf("error to create transfer")).MinTimes(1)
			},
			CheckResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, resp.Code)
				require.Contains(t, resp.Body.String(), "Service temporarily unavailable")
				require.Contains(t, resp.Body.String(), "error to create transfer")
			},
		},
	)

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			transferroute.Once = sync.Once{}
			m, server, ctrl := test.GetServerTest(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/transfers%s", transferroute.RootRoute)

			body, err := json.Marshal(tt.Body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			require.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			ctx := test.GetTestContext(t, req, recorder, true)

			if tt.SetupAuth != nil {
				tt.SetupAuth(ctx, t, req, m)
			}

			if tt.BuildMocks != nil {
				tt.BuildMocks(ctx, m, tt.Body)
			}

			server.Echo().ServeHTTP(recorder, req)
			if tt.CheckResponse != nil {
				tt.CheckResponse(t, recorder)
			}
		})
	}
}

func TestHandler_handleGetTransfers(t *testing.T) {
	tests := append(test.PrivateEndpointValidations,
		test.PrivateEndpointTest{
			Name: "Should pass with success",
			BuildMocks: func(ctx context.Context, m test.SvcMocks, _ any) {
				m.TransferAppMock.EXPECT().GetTransfers(ctx, int64(10), int64(0)).Return([]entity.Transfer{
					{TransferUUID: uuid.NewV4().String(), AccountOriginUUID: uuid.NewV4().String(), AccountDestinationUUID: uuid.NewV4().String(), Amount: 5.55, CreatedAt: time.Now()},
					{TransferUUID: uuid.NewV4().String(), AccountOriginUUID: uuid.NewV4().String(), AccountDestinationUUID: uuid.NewV4().String(), Amount: 7.77, CreatedAt: time.Now()},
				}, int64(0), nil).Times(1)
			},
			SetupAuth: func(ctx context.Context, t *testing.T, req *http.Request, m test.SvcMocks) {
				test.AddAuthorization(ctx, t, req, m)
			},
			CheckResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, resp.Code)
				require.NotEmpty(t, resp.Body)
				require.Contains(t, resp.Body.String(), "5.55")
				require.Contains(t, resp.Body.String(), "7.77")
			},
		},
		test.PrivateEndpointTest{
			Name: "Should return error if service get transfer returns error",
			BuildMocks: func(ctx context.Context, m test.SvcMocks, _ any) {
				m.TransferAppMock.EXPECT().GetTransfers(ctx, int64(10), int64(0)).Return(nil, int64(0), fmt.Errorf("error to get transfers")).Times(1)
			},
			SetupAuth: func(ctx context.Context, t *testing.T, req *http.Request, m test.SvcMocks) {
				test.AddAuthorization(ctx, t, req, m)
			},
			CheckResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, resp.Code)
				require.Contains(t, resp.Body.String(), "Service temporarily unavailable")
				require.Contains(t, resp.Body.String(), "error to get transfers")
			},
		},
	)
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			transferroute.Once = sync.Once{}
			m, server, ctrl := test.GetServerTest(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/transfers%s", transferroute.RootRoute)

			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			ctx := test.GetTestContext(t, req, recorder, true)

			if tt.SetupAuth != nil {
				tt.SetupAuth(ctx, t, req, m)
			}

			if tt.BuildMocks != nil {
				tt.BuildMocks(ctx, m, m)
			}

			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			server.Echo().ServeHTTP(recorder, req)

			if tt.CheckResponse != nil {
				tt.CheckResponse(t, recorder)
			}
		})
	}
}
