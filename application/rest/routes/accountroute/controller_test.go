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
	"github.com/diegoclair/go_boilerplate/mock"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

type mocks struct {
	mapper mapper.Mapper
	mas    *mock.MockAccountService
}

// var server *echo.Echo
// var accountMock mocks

func serverTest(t *testing.T) (accountMock mocks, server *echo.Echo) {

	ctrl := gomock.NewController(t)

	accountMock = mocks{
		mapper: mapper.New(),
		mas:    mock.NewMockAccountService(ctrl),
	}

	accountControler := NewController(accountMock.mas, accountMock.mapper)
	accountRoute := NewRouter(accountControler, "accounts")

	server = echo.New()
	appGroup := server.Group("/")

	accountRoute.RegisterRoutes(appGroup, nil)

	return accountMock, server
	//os.Exit(m.Run())
}

func TestController_handleAddAccount(t *testing.T) {

	type args struct {
		body any
	}
	tests := []struct {
		name          string
		args          args
		buildMocks    func(ctx context.Context, mock mocks, args args)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		// {
		// 	name: "Should complete request with no error",
		// 	args: args{
		// 		body: viewmodel.AddAccount{
		// 			Name: "Add withou Error",
		// 			Login: viewmodel.Login{
		// 				CPF:    "01234567890",
		// 				Secret: "secret@123",
		// 			},
		// 		},
		// 	},
		// 	buildMocks: func(ctx context.Context, mocks mocks, args args) {
		// 		body := args.body.(viewmodel.AddAccount)
		// 		mocks.mas.EXPECT().CreateAccount(ctx, entity.Account{Name: body.Name, CPF: body.CPF, Secret: body.Secret}).Times(1).Return(nil)
		// 	},
		// 	checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusCreated, resp.Code)
		// 		require.Empty(t, resp.Body)
		// 	},
		// },
		// {
		// 	name: "Should not be possible create an account without field name",
		// 	args: args{
		// 		body: viewmodel.AddAccount{
		// 			Name: "",
		// 			Login: viewmodel.Login{
		// 				CPF:    "01234567890",
		// 				Secret: "secret@123",
		// 			},
		// 		},
		// 	},
		// 	checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusUnprocessableEntity, resp.Code)
		// 		require.Contains(t, resp.Body.String(), "Invalid input data")
		// 		require.Contains(t, resp.Body.String(), "The field 'Name' is required")
		// 	},
		// },
		// {
		// 	name: "Should not be possible create an account without field CPF",
		// 	args: args{
		// 		body: viewmodel.AddAccount{
		// 			Name: "Teste name",
		// 			Login: viewmodel.Login{
		// 				CPF:    "",
		// 				Secret: "secret@123",
		// 			},
		// 		},
		// 	},
		// 	checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusUnprocessableEntity, resp.Code)
		// 		require.Contains(t, resp.Body.String(), "Invalid input data")
		// 		require.Contains(t, resp.Body.String(), "The field 'CPF' is required")
		// 	},
		// },
		// {
		// 	name: "Should not be possible create an account without field Secret",
		// 	args: args{
		// 		body: viewmodel.AddAccount{
		// 			Name: "Teste name",
		// 			Login: viewmodel.Login{
		// 				CPF:    "01234567890",
		// 				Secret: "",
		// 			},
		// 		},
		// 	},
		// 	checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusUnprocessableEntity, resp.Code)
		// 		require.Contains(t, resp.Body.String(), "Invalid input data")
		// 		require.Contains(t, resp.Body.String(), "The field 'Secret' is required")
		// 	},
		// },
		{
			name: "Should not be possible create an account with invalid cpf",
			args: args{
				body: viewmodel.AddAccount{
					Name: "Teste name",
					Login: viewmodel.Login{
						CPF:    "12345612345",
						Secret: "Secret@123",
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
						CPF:    "01234567890",
						Secret: "Secret@123",
					},
				},
			},
			buildMocks: func(ctx context.Context, mocks mocks, args args) {
				body := args.body.(viewmodel.AddAccount)
				mocks.mas.EXPECT().CreateAccount(ctx, entity.Account{Name: body.Name, CPF: body.CPF, Secret: body.Secret}).Times(1).Return(errors.New("some error"))
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, resp.Code)
				require.Contains(t, resp.Body.String(), "Service temporarily unavailable")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			accountMock, server := serverTest(t)

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
		buildMocks    func(ctx context.Context, mocks mocks, args args)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder, mock mocks, args args)
	}{
		{
			name: "Should complete request with no error",
			args: args{
				accountsToBuild: 2,
			},
			buildMocks: func(ctx context.Context, mock mocks, args args) {
				accounts := buildAccountsByQuantity(args.accountsToBuild)
				mock.mas.EXPECT().GetAccounts(ctx, int64(10), int64(0)).Times(1).Return(accounts, int64(2), nil)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder, mock mocks, args args) {
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			accountMock, server := serverTest(t)

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
