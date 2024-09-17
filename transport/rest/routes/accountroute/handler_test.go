package accountroute_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"

	"github.com/diegoclair/go_boilerplate/application/dto"
	"github.com/diegoclair/go_boilerplate/domain/entity"
	"github.com/diegoclair/go_boilerplate/transport/rest/routes/accountroute"
	"github.com/diegoclair/go_boilerplate/transport/rest/routes/test"
	"github.com/diegoclair/go_boilerplate/transport/rest/routeutils"
	"github.com/diegoclair/go_boilerplate/transport/rest/viewmodel"
	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestHandler_handleAddAccount(t *testing.T) {
	type args struct {
		body any
	}

	tests := []struct {
		name          string
		args          args
		buildMocks    func(ctx context.Context, m test.SvcMocks, args args)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Should complete request with no error",
			args: args{
				body: viewmodel.AddAccount{
					Name:     "Add without Error",
					CPF:      "01234567890",
					Password: "secret@123",
				},
			},
			buildMocks: func(ctx context.Context, m test.SvcMocks, args args) {
				body := args.body.(viewmodel.AddAccount)
				m.AccountAppMock.EXPECT().CreateAccount(ctx, body.ToDto()).Times(1).Return(nil)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, resp.Code)
				require.Empty(t, resp.Body)
			},
		},
		{
			name: "Should return error when body is invalid",
			args: args{
				body: "invalid body",
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, resp.Code)
				require.Contains(t, resp.Body.String(), "Unmarshal type error")
			},
		},
		{
			name: "Should return error if we have any error with service",
			args: args{
				body: viewmodel.AddAccount{
					Name:     "Error with service",
					CPF:      "01234567890",
					Password: "Secret@123",
				},
			},
			buildMocks: func(ctx context.Context, mock test.SvcMocks, args args) {
				body := args.body.(viewmodel.AddAccount)
				mock.AccountAppMock.EXPECT().CreateAccount(ctx, body.ToDto()).Times(1).Return(errors.New("some error"))
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, resp.Code)
				require.Contains(t, resp.Body.String(), "Service temporarily unavailable")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountroute.Once = sync.Once{}
			accountMock, server, ctrl := test.GetServerTest(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/accounts%s", accountroute.RootRoute)

			body, err := json.Marshal(tt.args.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			require.NoError(t, err)

			ctx := test.GetTestContext(t, req, recorder, false)

			if tt.buildMocks != nil {
				tt.buildMocks(ctx, accountMock, tt.args)
			}

			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			server.Echo().ServeHTTP(recorder, req)
			if tt.checkResponse != nil {
				tt.checkResponse(t, recorder)
			}
		})
	}
}

func buildAccountByID(id int) entity.Account {
	return entity.Account{UUID: "random", Name: "diego" + strconv.Itoa(id)}
}

func buildAccountsByQuantity(qtd int) (accounts []entity.Account) {
	for i := 0; i < qtd; i++ {
		accounts = append(accounts, buildAccountByID(i))
	}

	return
}

func TestHandler_GetAccounts(t *testing.T) {
	type args struct {
		page            int
		quantity        int
		accountsToBuild int
	}

	tests := []struct {
		name          string
		args          args
		buildMocks    func(ctx context.Context, mock test.SvcMocks, args args)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder, mock test.SvcMocks, args args)
	}{
		{
			name: "Should complete request with no error",
			args: args{
				accountsToBuild: 2,
			},
			buildMocks: func(ctx context.Context, mock test.SvcMocks, args args) {
				accounts := buildAccountsByQuantity(args.accountsToBuild)
				mock.AccountAppMock.EXPECT().GetAccounts(ctx, int64(10), int64(0)).Times(1).Return(accounts, int64(2), nil)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder, mock test.SvcMocks, args args) {
				require.Equal(t, http.StatusOK, resp.Code)
				accounts := buildAccountsByQuantity(args.accountsToBuild)
				take, skip := routeutils.GetTakeSkipFromPageQuantity(int64(args.page), int64(args.quantity))

				response := []viewmodel.AccountResponse{}
				for _, account := range accounts {
					item := viewmodel.AccountResponse{}
					item.FillFromEntity(account)
					response = append(response, item)
				}

				paginatedResp := viewmodel.BuildPaginatedResponse(response, skip, take, int64(args.accountsToBuild))
				expectedResp, err := json.Marshal(paginatedResp)
				require.NoError(t, err)
				require.Contains(t, resp.Body.String(), string(expectedResp))
			},
		},
		{
			name: "Should return error if we have some error with service",
			args: args{
				accountsToBuild: 0,
			},
			buildMocks: func(ctx context.Context, mock test.SvcMocks, args args) {
				accounts := buildAccountsByQuantity(args.accountsToBuild)
				mock.AccountAppMock.EXPECT().GetAccounts(ctx, int64(10), int64(0)).Times(1).Return(accounts, int64(0), fmt.Errorf("some service error"))
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder, mock test.SvcMocks, args args) {
				require.Equal(t, http.StatusServiceUnavailable, resp.Code)
				require.Contains(t, resp.Body.String(), "some service error")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountroute.Once = sync.Once{}
			accountMock, server, ctrl := test.GetServerTest(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/%s%s", accountroute.GroupRouteName, accountroute.RootRoute)

			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			ctx := test.GetTestContext(t, req, recorder, false)

			if tt.buildMocks != nil {
				tt.buildMocks(ctx, accountMock, tt.args)
			}

			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			server.Echo().ServeHTTP(recorder, req)
			if tt.checkResponse != nil {
				tt.checkResponse(t, recorder, accountMock, tt.args)
			}
		})
	}
}

func TestHandler_GetAccountByID(t *testing.T) {
	type args struct {
		accountUUID string
	}

	tests := []struct {
		name          string
		args          args
		buildMocks    func(ctx context.Context, mock test.SvcMocks, args args)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder, mock test.SvcMocks, args args)
	}{
		{
			name: "Should complete request with no error",
			args: args{
				accountUUID: "random",
			},
			buildMocks: func(ctx context.Context, mock test.SvcMocks, args args) {
				account := buildAccountByID(1)
				mock.AccountAppMock.EXPECT().GetAccountByUUID(ctx, args.accountUUID).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder, mock test.SvcMocks, args args) {
				require.Equal(t, http.StatusOK, resp.Code)
				account := buildAccountByID(1)

				response := viewmodel.AccountResponse{}
				response.FillFromEntity(account)

				expectedResp, err := json.Marshal(response)
				require.NoError(t, err)
				require.Contains(t, resp.Body.String(), string(expectedResp))
			},
		},
		{
			name: "Should return error if we have some error with service",
			args: args{
				accountUUID: "random",
			},
			buildMocks: func(ctx context.Context, mock test.SvcMocks, args args) {
				mock.AccountAppMock.EXPECT().GetAccountByUUID(ctx, args.accountUUID).Times(1).Return(entity.Account{}, fmt.Errorf("some service error"))
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder, mock test.SvcMocks, args args) {
				require.Equal(t, http.StatusServiceUnavailable, resp.Code)
				require.Contains(t, resp.Body.String(), "some service error")
			},
		},
		{
			name: "Should return error if we have an invalid uuid",
			args: args{accountUUID: " "},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder, mock test.SvcMocks, args args) {
				require.Equal(t, http.StatusUnprocessableEntity, resp.Code)
				require.Contains(t, resp.Body.String(), "Invalid account_uuid")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountroute.Once = sync.Once{}
			accountMock, server, ctrl := test.GetServerTest(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/%s/%s/", accountroute.GroupRouteName, tt.args.accountUUID)

			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			ctx := test.GetTestContext(t, req, recorder, false)

			if tt.buildMocks != nil {
				tt.buildMocks(ctx, accountMock, tt.args)
			}

			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			server.Echo().ServeHTTP(recorder, req)
			if tt.checkResponse != nil {
				tt.checkResponse(t, recorder, accountMock, tt.args)
			}
		})
	}
}

func TestHandler_handleAddBalance(t *testing.T) {
	type args struct {
		body        any
		accountUUID string
	}

	tests := []struct {
		name          string
		args          args
		buildMocks    func(ctx context.Context, mock test.SvcMocks, args args)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Should complete request with no error",
			args: args{
				body: viewmodel.AddBalance{
					Amount: 100,
				},
				accountUUID: "random",
			},
			buildMocks: func(ctx context.Context, mock test.SvcMocks, args args) {
				body := args.body.(viewmodel.AddBalance)

				mock.AccountAppMock.EXPECT().AddBalance(ctx, body.ToDto(args.accountUUID)).Times(1).Return(nil)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, resp.Code)
				require.Empty(t, resp.Body)
			},
		},
		{
			name: "Should return error when body is invalid",
			args: args{
				body: "invalid body",
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, resp.Code)
				require.Contains(t, resp.Body.String(), "Unmarshal type error")
			},
		},
		{
			name: "Should return error if we do not have an account_uuid in the url",
			args: args{
				body: viewmodel.AddBalance{
					Amount: 100,
				},
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnprocessableEntity, resp.Code)
				require.Contains(t, resp.Body.String(), "account_uuid is required")
			},
		},
		{
			name: "Should return error if we have any error with service",
			args: args{
				body: viewmodel.AddBalance{
					Amount: 100,
				},
				accountUUID: "random",
			},
			buildMocks: func(ctx context.Context, mock test.SvcMocks, args args) {
				input := dto.AddBalanceInput{
					AccountUUID: "random",
					Amount:      args.body.(viewmodel.AddBalance).Amount,
				}
				mock.AccountAppMock.EXPECT().AddBalance(ctx, input).Times(1).Return(errors.New("some error"))
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, resp.Code)
				require.Contains(t, resp.Body.String(), "Service temporarily unavailable")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountroute.Once = sync.Once{}
			accountMock, server, ctrl := test.GetServerTest(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/%s/%s/balance", accountroute.GroupRouteName, tt.args.accountUUID)

			body, err := json.Marshal(tt.args.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			require.NoError(t, err)

			ctx := test.GetTestContext(t, req, recorder, false)

			if tt.buildMocks != nil {
				tt.buildMocks(ctx, accountMock, tt.args)
			}

			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			server.Echo().ServeHTTP(recorder, req)
			if tt.checkResponse != nil {
				tt.checkResponse(t, recorder)
			}
		})
	}
}
