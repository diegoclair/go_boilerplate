package service

import (
	"log"
	"testing"

	"github.com/diegoclair/go-boilerplate/domain/contract"
	"github.com/diegoclair/go-boilerplate/infra/logger"
	"github.com/diegoclair/go-boilerplate/mock"
	"github.com/diegoclair/go-boilerplate/util/config"
	"github.com/golang/mock/gomock"
)

type serviceSetup struct {
	ctrl            *gomock.Controller
	log             logger.Logger
	dataManagerMock contract.DataManager
	mocks           mocks
}

type mocks struct {
	mar *mock.MockAccountRepo
	mcm *mock.MockCacheManager
}

func newServiceMock(t *testing.T) (mocks, *Service) {

	cfg, err := config.GetConfigEnvironment(config.ConfigDefaultFilepath)
	if err != nil {
		log.Fatalf("Error to load config: %v", err)
	}

	serviceSetup := serviceSetup{
		ctrl: gomock.NewController(t),
		log:  logger.New(*cfg),
	}

	serviceSetup.mocks = mocks{
		mar: mock.NewMockAccountRepo(serviceSetup.ctrl),
		mcm: mock.NewMockCacheManager(serviceSetup.ctrl),
	}

	serviceSetup.dataManagerMock = newDataMock(serviceSetup.ctrl, serviceSetup.mocks.mar)

	svc := New(serviceSetup.dataManagerMock, cfg, serviceSetup.mocks.mcm, serviceSetup.log)

	return serviceSetup.mocks, svc
}

type dataMock struct {
	ctrl   *gomock.Controller
	mar    *mock.MockAccountRepo
	mauthr *mock.MockAuthRepo
}

func newDataMock(ctrl *gomock.Controller, mar *mock.MockAccountRepo) contract.DataManager {
	return &dataMock{
		ctrl: ctrl,
		mar:  mar,
	}
}

func (d *dataMock) Begin() (contract.Transaction, error) {
	return mock.NewMockTransaction(d.ctrl), nil
}

func (d *dataMock) Account() contract.AccountRepo {
	return d.mar
}

func (d *dataMock) Auth() contract.AuthRepo {
	return d.mauthr
}
