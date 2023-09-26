package service

import (
	"testing"

	"github.com/diegoclair/go_boilerplate/infra/config"
	"github.com/diegoclair/go_boilerplate/infra/logger"
	"github.com/diegoclair/go_boilerplate/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type repoMock struct {
	mockDataManager *mocks.MockDataManager

	mockAuthRepo     *mocks.MockAuthRepo
	mockAccountRepo  *mocks.MockAccountRepo
	mockCacheManager *mocks.MockCacheManager
}

func newServiceTestMock(t *testing.T) (repoMocks repoMock, svc *service, ctrl *gomock.Controller) {

	cfg, err := config.GetConfigEnvironment("../../" + config.ConfigDefaultName)
	require.NoError(t, err)

	ctrl = gomock.NewController(t)
	log := logger.New(*cfg)

	repoMocks = repoMock{
		mockDataManager: mocks.NewMockDataManager(ctrl),

		mockAccountRepo:  mocks.NewMockAccountRepo(ctrl),
		mockCacheManager: mocks.NewMockCacheManager(ctrl),
		mockAuthRepo:     mocks.NewMockAuthRepo(ctrl),
	}

	repoMocks.mockDataManager.EXPECT().Account().Return(repoMocks.mockAccountRepo).AnyTimes()
	repoMocks.mockDataManager.EXPECT().Auth().Return(repoMocks.mockAuthRepo).AnyTimes()

	svc = newService(repoMocks.mockDataManager, cfg, repoMocks.mockCacheManager, log)

	return
}
