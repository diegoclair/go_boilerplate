package service

import (
	"testing"

	"github.com/diegoclair/go_boilerplate/domain/contract"
	"github.com/diegoclair/go_boilerplate/infra/logger"
	"github.com/diegoclair/go_boilerplate/mock"
	"github.com/diegoclair/go_boilerplate/util/config"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type mocks struct {
	mauthr *mock.MockAuthRepo
	mar    *mock.MockAccountRepo
	mcm    *mock.MockCacheManager
}

func newServiceTestMock(t *testing.T) (mocks, *Service) {

	cfg, err := config.GetConfigEnvironment("../../" + config.ConfigDefaultFilepath)
	require.NoError(t, err)

	ctrl := gomock.NewController(t)
	log := logger.New(*cfg)

	mocks := mocks{
		mar:    mock.NewMockAccountRepo(ctrl),
		mcm:    mock.NewMockCacheManager(ctrl),
		mauthr: mock.NewMockAuthRepo(ctrl),
	}

	dataManagerMock := newDataMock(ctrl, mocks)

	svc := New(dataManagerMock, cfg, mocks.mcm, log)

	return mocks, svc
}

type dataMock struct {
	ctrl  *gomock.Controller
	mocks mocks
}

func newDataMock(ctrl *gomock.Controller, mocks mocks) contract.DataManager {
	return &dataMock{
		ctrl:  ctrl,
		mocks: mocks,
	}
}

func (d *dataMock) Begin() (contract.Transaction, error) {
	return mock.NewMockTransaction(d.ctrl), nil
}

func (d *dataMock) Account() contract.AccountRepo {
	return d.mocks.mar
}

func (d *dataMock) Auth() contract.AuthRepo {
	return d.mocks.mauthr
}
