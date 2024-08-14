package pingroute

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/diegoclair/go_boilerplate/transport/rest/routeutils"
	"github.com/diegoclair/goswag"
	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestHandler_handlePing(t *testing.T) {
	tests := []struct {
		name          string
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "should return pong",
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				require.Equal(t, "{\"message\":\"pong\"}\n", recorder.Body.String())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/%s%s", RouteName, rootRoute)

			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			server := goswag.NewEcho()
			appGroup := server.Group("/")
			g := &routeutils.EchoGroups{
				AppGroup: appGroup,
			}

			pingroute := NewRouter(NewHandler(), RouteName)
			pingroute.RegisterRoutes(g)

			server.Echo().ServeHTTP(recorder, req)

			if tt.checkResponse != nil {
				tt.checkResponse(t, recorder)
			}
		})
	}
}
