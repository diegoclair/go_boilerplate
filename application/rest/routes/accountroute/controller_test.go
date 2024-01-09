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

	"github.com/IQ-tech/go-mapper"
	"github.com/diegoclair/go_boilerplate/application/rest/routeutils"
	"github.com/diegoclair/go_boilerplate/application/rest/viewmodel"
	"github.com/diegoclair/go_boilerplate/domain/entity"
	"github.com/diegoclair/go_boilerplate/mocks"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

type mock struct {
	mapper         mapper.Mapper
	accountService *mocks.MockAccountService
}

func getServerTest(t *testing.T) (accountMock mock, server *echo.Echo, ctrl *gomock.Controller) {

	ctrl = gomock.NewController(t)
	accountMock = mock{
		mapper:         mapper.New(),
		accountService: mocks.NewMockAccountService(ctrl),
	}

	transferControler := &Controller{accountMock.accountService, accountMock.mapper}
	transferRoute := NewRouter(transferControler, RouteName)

	server = echo.New()
	appGroup := server.Group("/")
	transferRoute.RegisterRoutes(appGroup, nil)
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
					Name: "Add withou Error",
					Login: viewmodel.Login{
						CPF:      "01234567890",
						Password: "secret@123",
					},
				},
			},
			buildMocks: func(ctx context.Context, mock mock, args args) {
				body := args.body.(viewmodel.AddAccount)
				mock.accountService.EXPECT().CreateAccount(ctx, entity.Account{Name: body.Name, CPF: body.CPF, Password: body.Password}).Times(1).Return(nil)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, resp.Code)
				require.Empty(t, resp.Body)
			},
		},
		{
			name: "Should not be possible create an account without field name",
			args: args{
				body: viewmodel.AddAccount{
					Name: "",
					Login: viewmodel.Login{
						CPF:      "01234567890",
						Password: "secret@123",
					},
				},
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnprocessableEntity, resp.Code)
				require.Contains(t, resp.Body.String(), "Invalid input data")
				require.Contains(t, resp.Body.String(), "The field 'Name' is required")
			},
		},
		{
			name: "Should not be possible create an account without field CPF",
			args: args{
				body: viewmodel.AddAccount{
					Name: "Teste name",
					Login: viewmodel.Login{
						CPF:      "",
						Password: "secret@123",
					},
				},
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnprocessableEntity, resp.Code)
				require.Contains(t, resp.Body.String(), "Invalid input data")
				require.Contains(t, resp.Body.String(), "The field 'CPF' is required")
			},
		},
		{
			name: "Should not be possible create an account without field Secret",
			args: args{
				body: viewmodel.AddAccount{
					Name: "Teste name",
					Login: viewmodel.Login{
						CPF:      "01234567890",
						Password: "",
					},
				},
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnprocessableEntity, resp.Code)
				require.Contains(t, resp.Body.String(), "Invalid input data")
				require.Contains(t, resp.Body.String(), "The field 'Password' is required")
			},
		},
		{
			name: "Should not be possible create an account with invalid cpf",
			args: args{
				body: viewmodel.AddAccount{
					Name: "Teste name",
					Login: viewmodel.Login{
						CPF:      "12345612345",
						Password: "Secret@123",
					},
				},
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnprocessableEntity, resp.Code)
				require.Contains(t, resp.Body.String(), "Invalid cpf document")
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
					Name: "Error with service",
					Login: viewmodel.Login{
						CPF:      "01234567890",
						Password: "Secret@123",
					},
				},
			},
			buildMocks: func(ctx context.Context, mock mock, args args) {
				body := args.body.(viewmodel.AddAccount)
				mock.accountService.EXPECT().CreateAccount(ctx, entity.Account{Name: body.Name, CPF: body.CPF, Password: body.Password}).Times(1).Return(errors.New("some error"))
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, resp.Code)
				require.Contains(t, resp.Body.String(), "Service temporarily unavailable")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

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
			server.ServeHTTP(recorder, req)
			if tt.checkResponse != nil {
				tt.checkResponse(t, recorder)
			}
		})
	}
}

func buildAccoununtByID(id int) entity.Account {
	return entity.Account{UUID: "random", Name: "diego" + strconv.Itoa(id)}
}
func buildAccountsByQuantity(qtd int) (accounts []entity.Account) {
	for i := 0; i < qtd; i++ {
		accounts = append(accounts, buildAccoununtByID(i))
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

				response := []viewmodel.Account{}
				err := mock.mapper.From(accounts).To(&response)
				require.NoError(t, err)

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
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder, mock mock, args args) {
				require.Equal(t, http.StatusServiceUnavailable, resp.Code)
				require.Contains(t, resp.Body.String(), "some service error")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			accountMock, server, ctrl := getServerTest(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/accounts%s", rootRoute)

			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			if tt.buildMocks != nil {
				e := echo.New()
				ctx := routeutils.GetContext(e.NewContext(req, recorder))
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
