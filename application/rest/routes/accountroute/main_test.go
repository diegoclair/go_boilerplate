package accountroute

import (
	"os"
	"testing"

	"github.com/IQ-tech/go-mapper"
	"github.com/diegoclair/go_boilerplate/mock"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
)

type mocks struct {
	mapper mapper.Mapper
	mas    *mock.MockAccountService
}

var server *echo.Echo
var accountMock mocks

func TestMain(m *testing.M) {

	ctrl := gomock.NewController(&testing.T{})

	accountMock = mocks{
		mapper: mapper.New(),
		mas:    mock.NewMockAccountService(ctrl),
	}

	accountControler := NewController(accountMock.mas, accountMock.mapper)
	accountRoute := NewRouter(accountControler, "accounts")

	server = echo.New()
	appGroup := server.Group("/")

	accountRoute.RegisterRoutes(appGroup, nil)

	os.Exit(m.Run())
}
