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
	"testing"

	"github.com/diegoclair/go_boilerplate/domain/account"
	"github.com/diegoclair/go_boilerplate/mocks"
	"github.com/diegoclair/go_boilerplate/transport/rest/routeutils"
	"github.com/diegoclair/go_boilerplate/transport/rest/viewmodel"
	"github.com/diegoclair/go_utils-lib/v2/validator"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

type mock struct {
	accountService *mocks.MockAccountService
}

func getServerTest(t *testing.T) (accountMock mock, server *echo.Echo, ctrl *gomock.Controller, accountController *Controller) {
	ctrl = gomock.NewController(t)
	accountMock = mock{
		accountService: mocks.NewMockAccountService(ctrl),
	}

	v, err := validator.NewValidator()
	require.NoError(t, err)

	accountController = &Controller{accountMock.accountService, routeutils.New(), v}
	accountRoute := NewRouter(accountController, RouteName)

	server = echo.New()
	appGroup := server.Group("/")
	g := &routeutils.EchoGroups{
		AppGroup: appGroup,
	}

	accountRoute.RegisterRoutes(g)
	return
}

func TestController_handleAddAccount(t *testing.T) {
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
				mock.accountService.EXPECT().CreateAccount(ctx, account.Account{Name: body.Name, CPF: body.CPF, Password: body.Password}).Times(1).Return(nil)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, resp.Code)
				require.Empty(t, resp.Body)
			},
		},
		{
			name: "Should validate Login required fields",
			args: args{
				body: viewmodel.AddAccount{},
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnprocessableEntity, resp.Code)
				require.Contains(t, resp.Body.String(), "Invalid input data")
				require.Contains(t, resp.Body.String(), "The field 'CPF' is required")
				require.Contains(t, resp.Body.String(), "The field 'Password' is required")
			},
		},
		{
			name: "Should validate required fields",
			args: args{
				body: viewmodel.AddAccount{
					CPF:      "01234567890",
					Password: "secret@123",
				},
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnprocessableEntity, resp.Code)
				require.Contains(t, resp.Body.String(), "Invalid input data")
				require.Contains(t, resp.Body.String(), "The field 'Name' is required")
			},
		},
		{
			name: "Should not be possible create an account with invalid cpf",
			args: args{
				body: viewmodel.AddAccount{
					Name:     "Teste name",
					CPF:      "12345612345",
					Password: "Secret@123",
				},
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnprocessableEntity, resp.Code)
				require.Contains(t, resp.Body.String(), "The field 'CPF' should be a valid cpf")
			},
		},
		{
			name: "Should return error with the body is empty",
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnprocessableEntity, resp.Code)
				require.Contains(t, resp.Body.String(), "Invalid input data")
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
				mock.accountService.EXPECT().CreateAccount(ctx, account.Account{Name: body.Name, CPF: body.CPF, Password: body.Password}).Times(1).Return(errors.New("some error"))
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, resp.Code)
				require.Contains(t, resp.Body.String(), "Service temporarily unavailable")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			accountMock, server, ctrl, s := getServerTest(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/accounts%s", rootRoute)

			body, err := json.Marshal(tt.args.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			require.NoError(t, err)

			if tt.buildMocks != nil {
				e := echo.New()
				ctx := s.utils.Req().GetContext(e.NewContext(req, recorder))
				tt.buildMocks(ctx, accountMock, tt.args)
			}

			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			server.ServeHTTP(recorder, req)
			if tt.checkResponse != nil {
				tt.checkResponse(t, recorder)
			}
		})
	}
}

func buildAccountByID(id int) account.Account {
	return account.Account{UUID: "random", Name: "diego" + strconv.Itoa(id)}
}

func buildAccountsByQuantity(qtd int) (accounts []account.Account) {
	for i := 0; i < qtd; i++ {
		accounts = append(accounts, buildAccountByID(i))
	}

	return
}

func TestController_GetAccounts(t *testing.T) {
	type args struct {
		page            int
		quantity        int
		accountsToBuild int
	}

	tests := []struct {
		name          string
		args          args
		buildMocks    func(ctx context.Context, mocks mock, args args)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder, mock mock, args args, s *Controller)
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
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder, mock mock, args args, s *Controller) {
				require.Equal(t, http.StatusOK, resp.Code)
				accounts := buildAccountsByQuantity(args.accountsToBuild)
				take, skip := s.utils.Req().GetTakeSkipFromPageQuantity(int64(args.page), int64(args.quantity))

				response := []viewmodel.Account{}
				for _, account := range accounts {
					item := viewmodel.Account{}
					item.FillFromEntity(account)
					response = append(response, item)
				}

				paginatedResp := routeutils.BuildPaginatedResult(response, skip, take, int64(args.accountsToBuild))
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
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder, mock mock, args args, s *Controller) {
				require.Equal(t, http.StatusServiceUnavailable, resp.Code)
				require.Contains(t, resp.Body.String(), "some service error")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			accountMock, server, ctrl, s := getServerTest(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/accounts%s", rootRoute)

			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			if tt.buildMocks != nil {
				e := echo.New()
				ctx := s.utils.Req().GetContext(e.NewContext(req, recorder))
				tt.buildMocks(ctx, accountMock, tt.args)
			}

			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			server.ServeHTTP(recorder, req)
			if tt.checkResponse != nil {
				tt.checkResponse(t, recorder, accountMock, tt.args, s)
			}
		})
	}
}

func TestController_GetAccountByID(t *testing.T) {
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

				response := viewmodel.Account{}
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
				mock.accountService.EXPECT().GetAccountByUUID(ctx, args.accountUUID).Times(1).Return(account.Account{}, fmt.Errorf("some service error"))
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

			accountMock, server, ctrl, s := getServerTest(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/%s/%s/", RouteName, tt.args.accountUUID)

			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			if tt.buildMocks != nil {
				e := echo.New()
				ctx := s.utils.Req().GetContext(e.NewContext(req, recorder))
				tt.buildMocks(ctx, accountMock, tt.args)
			}

			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			server.ServeHTTP(recorder, req)
			if tt.checkResponse != nil {
				tt.checkResponse(t, recorder, accountMock, tt.args)
			}
		})
	}
}
