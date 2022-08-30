package testutil

import (
	"log"
	"testing"
	"time"

	"github.com/diegoclair/go_boilerplate/infra/auth"
	"github.com/diegoclair/go_boilerplate/mock"
	"github.com/diegoclair/go_boilerplate/util/config"
	"github.com/golang/mock/gomock"
)

type Mocks struct {
	AccountServiceMock  *mock.MockAccountService
	AuthServiceMock     *mock.MockAuthService
	TransferServiceMock *mock.MockTransferService
}

func NewServiceManagerTest(t *testing.T) Mocks {

	ctrl := gomock.NewController(t)

	mocks := Mocks{
		AccountServiceMock:  mock.NewMockAccountService(ctrl),
		AuthServiceMock:     mock.NewMockAuthService(ctrl),
		TransferServiceMock: mock.NewMockTransferService(ctrl),
	}

	return mocks
}

func GetTokenMaker() auth.AuthToken {
	cfg, err := config.GetConfigEnvironment("../../../../" + config.ConfigDefaultFilepath)
	if err != nil {
		log.Fatal("failed to get config: ", err)
	}

	cfg.App.Auth.AccessTokenDuration = 2 * time.Second
	cfg.App.Auth.RefreshTokenDuration = 2 * time.Second

	tokenMaker, err := auth.NewAuthToken(cfg.App.Auth)
	if err != nil {
		log.Fatal("failed to create authToken: ", err)
	}

	return tokenMaker
}
