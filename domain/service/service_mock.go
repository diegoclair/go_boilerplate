package service

import (
	"testing"

	"github.com/diegoclair/go_boilerplate/domain/contract"
	"github.com/diegoclair/go_boilerplate/infra/config"
	"github.com/diegoclair/go_boilerplate/infra/logger"
	"github.com/diegoclair/go_boilerplate/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type repoMock struct {
	mockAuthRepo     *mocks.MockAuthRepo
	mockAccountRepo  *mocks.MockAccountRepo
	mockCacheManager *mocks.MockCacheManager
}

func newServiceTestMock(t *testing.T) (repoMocks repoMock, svc *Service, ctrl *gomock.Controller) {

	cfg, err := config.GetConfigEnvironment("../../" + config.ConfigDefaultName)
	require.NoError(t, err)

	ctrl = gomock.NewController(t)
	log := logger.New(*cfg)

	repoMocks = repoMock{
		mockAccountRepo:  mocks.NewMockAccountRepo(ctrl),
		mockCacheManager: mocks.NewMockCacheManager(ctrl),
		mockAuthRepo:     mocks.NewMockAuthRepo(ctrl),
	}

	dataManagerMock := newDataMock(repoMocks)

	svc = New(dataManagerMock, cfg, repoMocks.mockCacheManager, log)

	return
}

type dataMock struct {
	mocks repoMock
}

func newDataMockTransaction(mocks repoMock) contract.Transaction {
	return &dataMock{

		mocks: mocks,
	}
}

func (d *dataMock) Rollback() error {
	return nil
}

func (d *dataMock) Commit() error {
	return nil
}

func newDataMock(mocks repoMock) contract.DataManager {
	return &dataMock{
		mocks: mocks,
	}
}

func (d *dataMock) Begin() (contract.Transaction, error) {
	return newDataMockTransaction(d.mocks), nil
}

func (d *dataMock) Account() contract.AccountRepo {
	return d.mocks.mockAccountRepo
}

func (d *dataMock) Auth() contract.AuthRepo {
	return d.mocks.mockAuthRepo
}
