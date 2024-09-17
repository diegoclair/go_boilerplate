package servermiddleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/diegoclair/go_boilerplate/infra"
	"github.com/diegoclair/go_boilerplate/infra/contract"
	infraMocks "github.com/diegoclair/go_boilerplate/infra/mocks"
	"github.com/diegoclair/go_utils/resterrors"
	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAuthMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAuthToken := infraMocks.NewMockAuthToken(ctrl)
	cacheMock := infraMocks.NewMockCacheManager(ctrl)
	middleware := AuthMiddlewarePrivateRoute(mockAuthToken, cacheMock)

	t.Run("Should complete the middleware without errors", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(infra.TokenKey.String(), "Bearer")
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		mockAuthToken.EXPECT().VerifyToken(gomock.Any(), "Bearer").Return(&contract.TokenPayload{
			AccountUUID: "uuid",
			SessionUUID: "session",
		}, nil)

		cacheMock.EXPECT().GetString(gomock.Any(), "Bearer").Return("", nil)
		err := middleware(func(c echo.Context) error {
			return nil
		})(c)

		assert.Nil(t, err)
		assert.Equal(t, "uuid", c.Get(infra.AccountUUIDKey.String()))
		assert.Equal(t, "session", c.Get(infra.SessionKey.String()))
	})
	t.Run("Should return error when access token is required", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := middleware(func(c echo.Context) error {
			return nil
		})(c)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusUnauthorized, err.(resterrors.RestErr).StatusCode())
	})

	t.Run("Should return error when verify token fails", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(infra.TokenKey.String(), "Bearer")
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		mockAuthToken.EXPECT().VerifyToken(gomock.Any(), "Bearer").Return(nil, resterrors.NewUnauthorizedError("an error occurred"))

		err := middleware(func(c echo.Context) error {
			return nil
		})(c)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusUnauthorized, err.(resterrors.RestErr).StatusCode())
	})

	t.Run("Should return error when token is already invalid", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(infra.TokenKey.String(), "Bearer")
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		mockAuthToken.EXPECT().VerifyToken(gomock.Any(), "Bearer").Return(&contract.TokenPayload{
			AccountUUID: "uuid",
			SessionUUID: "session",
		}, nil)

		cacheMock.EXPECT().GetString(gomock.Any(), "Bearer").Return("invalid", nil)
		err := middleware(func(c echo.Context) error {
			return nil
		})(c)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusUnauthorized, err.(resterrors.RestErr).StatusCode())
	})
}
