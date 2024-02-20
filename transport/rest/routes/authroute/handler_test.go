package authroute

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/diegoclair/go_boilerplate/application/dto"
	"github.com/diegoclair/go_boilerplate/domain/account"
	"github.com/diegoclair/go_boilerplate/infra/auth"
	"github.com/diegoclair/go_boilerplate/mocks"
	"github.com/diegoclair/go_boilerplate/transport/rest/routeutils"
	"github.com/diegoclair/go_boilerplate/transport/rest/viewmodel"
	"github.com/diegoclair/go_utils/validator"
	"github.com/diegoclair/goswag"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

type mock struct {
	authService *mocks.MockAuthService
	authToken   *mocks.MockAuthToken
}

func getServerTest(t *testing.T) (authMock mock, server goswag.Echo, ctrl *gomock.Controller, accountHandler *Handler) {
	ctrl = gomock.NewController(t)
	authMock = mock{
		authService: mocks.NewMockAuthService(ctrl),
		authToken:   mocks.NewMockAuthToken(ctrl),
	}

	v, err := validator.NewValidator()
	require.NoError(t, err)

	accountHandler = &Handler{authMock.authService, authMock.authToken, routeutils.New(), v}
	accountRoute := NewRouter(accountHandler, RouteName)

	server = goswag.NewEcho()
	appGroup := server.Group("/")
	g := &routeutils.EchoGroups{
		AppGroup: appGroup,
	}

	accountRoute.RegisterRoutes(g)
	return
}

func TestHandler_handleLogin(t *testing.T) {
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
				body: viewmodel.Login{
					CPF:      "01234567890",
					Password: "12345678",
				},
			},
			buildMocks: func(ctx context.Context, mock mock, args args) {
				input := args.body.(viewmodel.Login)
				mock.authService.EXPECT().Login(ctx, input.CPF, input.Password).Return(account.Account{ID: 1, UUID: "uuid"}, nil).Times(1)
				mock.authToken.EXPECT().CreateAccessToken(ctx, "uuid", gomock.Any()).Return("a123", &auth.TokenPayload{}, nil).Times(1)
				mock.authToken.EXPECT().CreateRefreshToken(ctx, "uuid", gomock.Any()).Return("r123", &auth.TokenPayload{ExpiredAt: time.Now()}, nil).Times(1)
				mock.authService.EXPECT().CreateSession(ctx, gomock.Any()).DoAndReturn(
					func(ctx context.Context, req dto.Session) error {
						require.NotEmpty(t, req.SessionUUID)
						require.Equal(t, int64(1), req.AccountID)
						require.Equal(t, "r123", req.RefreshToken)
						require.NotEmpty(t, req.RefreshTokenExpiredAt)

						return nil
					},
				).Times(1)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authMock, server, ctrl, authHandler := getServerTest(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/%s%s", RouteName, loginRoute)

			body, err := json.Marshal(tt.args.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			require.NoError(t, err)

			if tt.buildMocks != nil {
				e := echo.New()
				ctx := authHandler.utils.Req().GetContext(e.NewContext(req, recorder))
				tt.buildMocks(ctx, authMock, tt.args)
			}

			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			server.Echo().ServeHTTP(recorder, req)
			if tt.checkResponse != nil {
				tt.checkResponse(t, recorder)
			}
		})
	}
}
