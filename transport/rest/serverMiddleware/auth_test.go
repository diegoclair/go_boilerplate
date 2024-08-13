package servermiddleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/diegoclair/go_boilerplate/infra"
	"github.com/diegoclair/go_boilerplate/infra/auth"
	"github.com/diegoclair/go_boilerplate/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAuthMiddleware(t *testing.T) {
	mockAuthToken := mocks.NewMockAuthToken(gomock.NewController(t))
	middleware := AuthMiddlewarePrivateRoute(mockAuthToken)

	t.Run("Should return error when access token is required", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := middleware(func(c echo.Context) error {
			return nil
		})(c)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusUnauthorized, err.(*echo.HTTPError).Code)
	})

	t.Run("Should return error when verify token fails", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(infra.TokenKey.String(), "Bearer")
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		mockAuthToken.EXPECT().VerifyToken(gomock.Any(), "Bearer").Return(nil, assert.AnError)

		err := middleware(func(c echo.Context) error {
			return nil
		})(c)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusUnauthorized, err.(*echo.HTTPError).Code)
	})

	t.Run("Should complete the middleware without errors", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(infra.TokenKey.String(), "Bearer")
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		mockAuthToken.EXPECT().VerifyToken(gomock.Any(), "Bearer").Return(&auth.TokenPayload{
			AccountUUID: "uuid",
			SessionUUID: "session",
		}, nil)

		err := middleware(func(c echo.Context) error {
			return nil
		})(c)

		assert.Nil(t, err)
		assert.Equal(t, "uuid", c.Get(infra.AccountUUIDKey.String()))
		assert.Equal(t, "session", c.Get(infra.SessionKey.String()))
	})
}
