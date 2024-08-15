package accountroute

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
	"github.com/diegoclair/go_boilerplate/mocks"
	"github.com/diegoclair/go_boilerplate/transport/rest/routeutils"
	"github.com/diegoclair/go_boilerplate/transport/rest/viewmodel"
	"github.com/diegoclair/goswag"
	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type mock struct {
	accountService *mocks.MockAccountService
}

func getServerTest(t *testing.T) (accountMock mock, server goswag.Echo, ctrl *gomock.Controller) {
	ctrl = gomock.NewController(t)
	accountMock = mock{
		accountService: mocks.NewMockAccountService(ctrl),
	}

	accountHandler := NewHandler(accountMock.accountService)
	accountRoute := NewRouter(accountHandler, RouteName)

	server = goswag.NewEcho()
	appGroup := server.Group("/")
	g := &routeutils.EchoGroups{
		AppGroup: appGroup,
	}

	accountRoute.RegisterRoutes(g)
	return
}

func TestHandler_handleAddAccount(t *testing.T) {
	type args struct {
		body any
	}

	tests := []struct {
		name          string
		args          args
		buildMocks    func(ctx context.Context, mock mock, args args)
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
			buildMocks: func(ctx context.Context, mock mock, args args) {
				body := args.body.(viewmodel.AddAccount)
				mock.accountService.EXPECT().CreateAccount(ctx, body.ToDto()).Times(1).Return(nil)
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
			buildMocks: func(ctx context.Context, mock mock, args args) {
				body := args.body.(viewmodel.AddAccount)
				mock.accountService.EXPECT().CreateAccount(ctx, body.ToDto()).Times(1).Return(errors.New("some error"))
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, resp.Code)
				require.Contains(t, resp.Body.String(), "Service temporarily unavailable")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			once = sync.Once{}
			accountMock, server, ctrl := getServerTest(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/accounts%s", rootRoute)

			body, err := json.Marshal(tt.args.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			require.NoError(t, err)

			if tt.buildMocks != nil {
				e := echo.New()
				ctx := routeutils.GetContext(e.NewContext(req, recorder))
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
		buildMocks    func(ctx context.Context, mocks mock, args args)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder, mock mock, args args)
	}{
		{
			name: "Should complete request with no error",
			args: args{
				accountsToBuild: 2,
			},
			buildMocks: func(ctx context.Context, mock mock, args args) {
				accounts := buildAccountsByQuantity(args.accountsToBuild)
				mock.accountService.EXPECT().GetAccounts(ctx, int64(10), int64(0)).Times(1).Return(accounts, int64(2), nil)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder, mock mock, args args) {
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
			buildMocks: func(ctx context.Context, mock mock, args args) {
				accounts := buildAccountsByQuantity(args.accountsToBuild)
				mock.accountService.EXPECT().GetAccounts(ctx, int64(10), int64(0)).Times(1).Return(accounts, int64(0), fmt.Errorf("some service error"))
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder, mock mock, args args) {
				require.Equal(t, http.StatusServiceUnavailable, resp.Code)
				require.Contains(t, resp.Body.String(), "some service error")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			once = sync.Once{}
			accountMock, server, ctrl := getServerTest(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/%s%s", RouteName, rootRoute)

			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			if tt.buildMocks != nil {
				e := echo.New()
				ctx := routeutils.GetContext(e.NewContext(req, recorder))
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
		buildMocks    func(ctx context.Context, mocks mock, args args)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder, mock mock, args args)
	}{
		{
			name: "Should complete request with no error",
			args: args{
				accountUUID: "random",
			},
			buildMocks: func(ctx context.Context, mock mock, args args) {
				account := buildAccountByID(1)
				mock.accountService.EXPECT().GetAccountByUUID(ctx, args.accountUUID).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder, mock mock, args args) {
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
			buildMocks: func(ctx context.Context, mock mock, args args) {
				mock.accountService.EXPECT().GetAccountByUUID(ctx, args.accountUUID).Times(1).Return(entity.Account{}, fmt.Errorf("some service error"))
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder, mock mock, args args) {
				require.Equal(t, http.StatusServiceUnavailable, resp.Code)
				require.Contains(t, resp.Body.String(), "some service error")
			},
		},
		{
			name: "Should return error if we have an invalid uuid",
			args: args{accountUUID: " "},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder, mock mock, args args) {
				require.Equal(t, http.StatusUnprocessableEntity, resp.Code)
				require.Contains(t, resp.Body.String(), "Invalid account_uuid")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			once = sync.Once{}
			accountMock, server, ctrl := getServerTest(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/%s/%s/", RouteName, tt.args.accountUUID)

			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			if tt.buildMocks != nil {
				e := echo.New()
				ctx := routeutils.GetContext(e.NewContext(req, recorder))
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
		buildMocks    func(ctx context.Context, mock mock, args args)
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
			buildMocks: func(ctx context.Context, mock mock, args args) {
				body := args.body.(viewmodel.AddBalance)

				mock.accountService.EXPECT().AddBalance(ctx, body.ToDto(args.accountUUID)).Times(1).Return(nil)
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
			buildMocks: func(ctx context.Context, mock mock, args args) {
				input := dto.AddBalanceInput{
					AccountUUID: "random",
					Amount:      args.body.(viewmodel.AddBalance).Amount,
				}
				mock.accountService.EXPECT().AddBalance(ctx, input).Times(1).Return(errors.New("some error"))
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, resp.Code)
				require.Contains(t, resp.Body.String(), "Service temporarily unavailable")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			once = sync.Once{}
			accountMock, server, ctrl := getServerTest(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/%s/%s/balance", RouteName, tt.args.accountUUID)

			body, err := json.Marshal(tt.args.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			require.NoError(t, err)

			if tt.buildMocks != nil {
				e := echo.New()
				ctx := routeutils.GetContext(e.NewContext(req, recorder))
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
