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
	"github.com/diegoclair/go_boilerplate/transport/rest/routes/shared"
	"github.com/diegoclair/go_boilerplate/transport/rest/routes/transferroute"
	"github.com/diegoclair/go_boilerplate/transport/rest/viewmodel"
	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/twinj/uuid"
)

func TestHandler_handleAddTransfer(t *testing.T) {
	type args struct {
		body any
	}

	tests := []struct {
		name          string
		args          args
		setupAuth     func(ctx context.Context, t *testing.T, req *http.Request)
		buildMocks    func(ctx context.Context, m shared.SvcMocks, args args)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Should complete request with no error",
			args: args{
				body: viewmodel.TransferReq{
					AccountDestinationUUID: "randomUUID",
					Amount:                 5.55,
				},
			},
			setupAuth: func(ctx context.Context, t *testing.T, req *http.Request) {
				shared.AddAuthorization(ctx, t, req)
			},
			buildMocks: func(ctx context.Context, m shared.SvcMocks, args args) {
				body := args.body.(viewmodel.TransferReq)
				m.TransferMock.EXPECT().CreateTransfer(ctx,
					dto.TransferInput{AccountDestinationUUID: body.AccountDestinationUUID, Amount: body.Amount}).
					Return(nil).MinTimes(1)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, resp.Code)
				require.Empty(t, resp.Body)
			},
		},
		{
			name: "Should return error if body is invalid",
			args: args{
				body: "invalid body",
			},
			setupAuth: func(ctx context.Context, t *testing.T, req *http.Request) {
				shared.AddAuthorization(ctx, t, req)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, resp.Code)
				require.Contains(t, resp.Body.String(), "Unmarshal type error")
			},
		},
		{
			name: "Should return error id we have some error on create transfer",
			args: args{
				body: viewmodel.TransferReq{
					AccountDestinationUUID: "randomUUID2",
					Amount:                 8.88,
				},
			},
			setupAuth: func(ctx context.Context, t *testing.T, req *http.Request) {
				shared.AddAuthorization(ctx, t, req)
			},
			buildMocks: func(ctx context.Context, m shared.SvcMocks, args args) {
				body := args.body.(viewmodel.TransferReq)
				m.TransferMock.EXPECT().CreateTransfer(ctx,
					dto.TransferInput{AccountDestinationUUID: body.AccountDestinationUUID, Amount: body.Amount}).
					Return(fmt.Errorf("error to create transfer")).MinTimes(1)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, resp.Code)
				require.Contains(t, resp.Body.String(), "Service temporarily unavailable")
				require.Contains(t, resp.Body.String(), "error to create transfer")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			transferroute.Once = sync.Once{}
			m, server, ctrl := shared.GetServerTest(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/transfers%s", transferroute.RootRoute)

			body, err := json.Marshal(tt.args.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			require.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			ctx := shared.GetTestContext(t, req, recorder)

			if tt.setupAuth != nil {
				tt.setupAuth(ctx, t, req)
			}

			if tt.buildMocks != nil {
				tt.buildMocks(ctx, m, tt.args)
			}

			server.Echo().ServeHTTP(recorder, req)
			if tt.checkResponse != nil {
				tt.checkResponse(t, recorder)
			}
		})
	}
}

func TestHandler_handleGetTransfers(t *testing.T) {
	tests := []struct {
		name          string
		buildMocks    func(ctx context.Context, m shared.SvcMocks)
		setupAuth     func(ctx context.Context, t *testing.T, req *http.Request)
		checkResponse func(t *testing.T, resp *httptest.ResponseRecorder)
		sleep         bool
	}{
		{
			name: "Should pass with success",
			buildMocks: func(ctx context.Context, m shared.SvcMocks) {
				m.TransferMock.EXPECT().GetTransfers(ctx, int64(10), int64(0)).Return([]entity.Transfer{
					{TransferUUID: uuid.NewV4().String(), AccountOriginUUID: uuid.NewV4().String(), AccountDestinationUUID: uuid.NewV4().String(), Amount: 5.55, CreatedAt: time.Now()},
					{TransferUUID: uuid.NewV4().String(), AccountOriginUUID: uuid.NewV4().String(), AccountDestinationUUID: uuid.NewV4().String(), Amount: 7.77, CreatedAt: time.Now()},
				}, int64(0), nil).Times(1)
			},
			setupAuth: func(ctx context.Context, t *testing.T, req *http.Request) {
				shared.AddAuthorization(ctx, t, req)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, resp.Code)
				require.NotEmpty(t, resp.Body)
				require.Contains(t, resp.Body.String(), "5.55")
				require.Contains(t, resp.Body.String(), "7.77")
			},
		},
		{
			name: "Should return expired token error",
			setupAuth: func(ctx context.Context, t *testing.T, req *http.Request) {
				shared.AddAuthorization(ctx, t, req)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, resp.Code)
				require.Contains(t, resp.Body.String(), "token has expired")
			},
			sleep: true,
		},
		{
			name: "Should return error if service get transfer returns error",
			buildMocks: func(ctx context.Context, m shared.SvcMocks) {
				m.TransferMock.EXPECT().GetTransfers(ctx, int64(10), int64(0)).Return(nil, int64(0), fmt.Errorf("error to get transfers")).Times(1)
			},
			setupAuth: func(ctx context.Context, t *testing.T, req *http.Request) {
				shared.AddAuthorization(ctx, t, req)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, resp.Code)
				require.Contains(t, resp.Body.String(), "Service temporarily unavailable")
				require.Contains(t, resp.Body.String(), "error to get transfers")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			transferroute.Once = sync.Once{}
			m, server, ctrl := shared.GetServerTest(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/transfers%s", transferroute.RootRoute)

			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			ctx := shared.GetTestContext(t, req, recorder)

			if tt.setupAuth != nil {
				tt.setupAuth(ctx, t, req)
			}

			if tt.buildMocks != nil {
				tt.buildMocks(ctx, m)
			}

			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			if tt.sleep {
				time.Sleep(2 * time.Second)
			}
			server.Echo().ServeHTTP(recorder, req)

			if tt.checkResponse != nil {
				tt.checkResponse(t, recorder)
			}
		})
	}
}
